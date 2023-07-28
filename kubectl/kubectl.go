package kubectl

import (
	log "github.com/sirupsen/logrus"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/apply"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

func ApplyYamlFile(fileName string) {

	kubeConfigFlags := genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag()
	matchVersionKubeConfigFlags := cmdutil.NewMatchVersionFlags(kubeConfigFlags)
	f := cmdutil.NewFactory(matchVersionKubeConfigFlags)

	IOStreams, _, out, _ := genericclioptions.NewTestIOStreams()

	cmd := apply.NewCmdApply("kubectl", f, IOStreams)
	cmd.Flags().Set("filename", fileName)
	cmd.Run(cmd, []string{})

	log.Info(out.String())
	out.Reset()
}
