package testlib

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"path"
	goruntime "runtime"

	"github.com/vmware/govmomi/vapi/rest"
	vapitags "github.com/vmware/govmomi/vapi/tags"

	ocpv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/vsphere-problem-detector/pkg/util"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/session"
	"github.com/vmware/govmomi/simulator"

	// required to initialize the VAPI endpoint.
	_ "github.com/vmware/govmomi/vapi/simulator"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	"gopkg.in/gcfg.v1"
	v1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/legacy-cloud-providers/vsphere"
)

const (
	DefaultModel    = "testlib/testdata/default"
	defaultDC       = "DC0"
	defaultVMPath   = "/DC0/vm/"
	defaultHost     = "H0"
	DefaultHostId   = "host-24" // Generated by vcsim
	defaultHostPath = "/DC0/host/DC0_"
)

var (
	NodeProperties = []string{"config.extraConfig", "config.flags", "config.version", "runtime.host"}
)

type SimulatedVM struct {
	Name, UUID string
}

var (
	// Virtual machines generated by vSphere simulator. UUIDs look generated, but they're stable.
	defaultVMs = []SimulatedVM{
		{"DC0_H0_VM0", "265104de-1472-547c-b873-6dc7883fb6cb"},
		{"DC0_H0_VM1", "12f8928d-f144-5c57-89db-dd2d0902c9fa"},
	}
)

func init() {
	_, filename, _, _ := goruntime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func connectToSimulator(s *simulator.Server) (*vim25.Client, error) {
	client, err := govmomi.NewClient(context.TODO(), s.URL, true)
	if err != nil {
		return nil, err
	}
	return client.Client, nil
}

func simulatorConfig() *vsphere.VSphereConfig {
	var cfg vsphere.VSphereConfig
	// Configuration that corresponds to the simulated vSphere
	data := `[Global]
secret-name = "vsphere-creds"
secret-namespace = "kube-system"
insecure-flag = "1"

[Workspace]
server = "localhost"
datacenter = "DC0"
default-datastore = "LocalDS_0"
folder = "/DC0/vm"
resourcepool-path = "/DC0/host/DC0_H0/Resources"

[VirtualCenter "dc0"]
datacenters = "DC0"
`
	err := gcfg.ReadStringInto(&cfg, data)
	if err != nil {
		panic(err)
	}
	return &cfg
}

type TestSetup struct {
	Context     context.Context
	Username    string
	VMConfig    *vsphere.VSphereConfig
	VMClient    *vim25.Client
	TagManager  *vapitags.Manager
	ClusterInfo *util.ClusterInfo
}

func SetupSimulator(kubeClient *FakeKubeClient, modelDir string) (setup *TestSetup, cleanup func(), err error) {
	model := simulator.Model{}
	err = model.Load(modelDir)
	if err != nil {
		return nil, nil, err
	}
	model.Service.TLS = new(tls.Config)
	model.Service.RegisterEndpoints = true

	s := model.Service.NewServer()
	client, err := connectToSimulator(s)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to the simulator: %s", err)
	}
	clusterInfo := util.NewClusterInfo()

	sessionMgr := session.NewManager(client)
	userSession, err := sessionMgr.UserSession(context.TODO())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to setup user session: %v", err)
	}

	restClient := rest.NewClient(client)
	restClient.Login(context.TODO(), s.URL.User)
	testSetup := &TestSetup{
		Context:     context.TODO(),
		VMConfig:    simulatorConfig(),
		TagManager:  vapitags.NewManager(restClient),
		ClusterInfo: clusterInfo,
		VMClient:    client,
	}
	testSetup.VMConfig.Workspace.VCenterIP = "dc0"
	testSetup.VMConfig.VirtualCenter["dc0"].User = userSession.UserName
	testSetup.Username = userSession.UserName

	cleanup = func() {
		s.Close()
		model.Remove()
	}
	return testSetup, cleanup, nil
}

type FakeKubeClient struct {
	Infrastructure *ocpv1.Infrastructure
	Nodes          []*v1.Node
	StorageClasses []*storagev1.StorageClass
	PVs            []*v1.PersistentVolume
}

func (f *FakeKubeClient) GetInfrastructure(ctx context.Context) (*ocpv1.Infrastructure, error) {
	return f.Infrastructure, nil
}

func (f *FakeKubeClient) ListNodes(ctx context.Context) ([]*v1.Node, error) {
	return f.Nodes, nil
}

func (f *FakeKubeClient) ListStorageClasses(ctx context.Context) ([]*storagev1.StorageClass, error) {
	return f.StorageClasses, nil
}

func (f *FakeKubeClient) ListPVs(ctx context.Context) ([]*v1.PersistentVolume, error) {
	return f.PVs, nil
}

func Node(name string, modifiers ...func(*v1.Node)) *v1.Node {
	n := &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: v1.NodeSpec{
			ProviderID: "",
		},
	}
	for _, modifier := range modifiers {
		modifier(n)
	}
	return n
}

func WithProviderID(id string) func(*v1.Node) {
	return func(node *v1.Node) {
		node.Spec.ProviderID = id
	}
}

func DefaultNodes() []*v1.Node {
	nodes := []*v1.Node{}
	for _, vm := range defaultVMs {
		node := Node(vm.Name, WithProviderID("vsphere://"+vm.UUID))
		nodes = append(nodes, node)
	}
	return nodes
}

func Infrastructure(modifiers ...func(*ocpv1.Infrastructure)) *ocpv1.Infrastructure {
	infra := &ocpv1.Infrastructure{
		ObjectMeta: metav1.ObjectMeta{
			Name: "cluster",
		},
		Spec: ocpv1.InfrastructureSpec{
			PlatformSpec: ocpv1.PlatformSpec{
				VSphere: &ocpv1.VSpherePlatformSpec{},
			},
		},
		Status: ocpv1.InfrastructureStatus{
			InfrastructureName: "my-cluster-id",
		},
	}

	for _, modifier := range modifiers {
		modifier(infra)
	}
	return infra
}

func GetVM(vmClient *vim25.Client, node *v1.Node) (*mo.VirtualMachine, error) {
	finder := find.NewFinder(vmClient, true)
	vm, err := finder.VirtualMachine(context.TODO(), defaultVMPath+node.Name)
	if err != nil {
		return nil, err
	}

	var o mo.VirtualMachine
	err = vm.Properties(context.TODO(), vm.Reference(), NodeProperties, &o)
	if err != nil {
		return nil, fmt.Errorf("failed to load VM %s: %s", node.Name, err)
	}

	return &o, nil
}

func CustomizeVM(vmClient *vim25.Client, node *v1.Node, spec *types.VirtualMachineConfigSpec) error {
	finder := find.NewFinder(vmClient, true)
	vm, err := finder.VirtualMachine(context.TODO(), defaultVMPath+node.Name)
	if err != nil {
		return err
	}

	task, err := vm.Reconfigure(context.TODO(), *spec)
	if err != nil {
		return err
	}

	err = task.Wait(context.TODO())
	return err
}

func SetHardwareVersion(vmClient *vim25.Client, node *v1.Node, hardwareVersion string) error {
	err := CustomizeVM(vmClient, node, &types.VirtualMachineConfigSpec{
		ExtraConfig: []types.BaseOptionValue{
			&types.OptionValue{
				Key: "SET.config.version", Value: hardwareVersion,
			},
		}})
	return err
}

func CustomizeHostVersion(hostSystemId string, version string, apiVersion string) error {
	hsRef := simulator.Map.Get(types.ManagedObjectReference{
		Type:  "HostSystem",
		Value: hostSystemId,
	})
	if hsRef == nil {
		return fmt.Errorf("can't find HostSystem %s", hostSystemId)
	}

	hs := hsRef.(*simulator.HostSystem)
	hs.Config.Product.Version = version
	hs.Config.Product.ApiVersion = apiVersion
	return nil
}
