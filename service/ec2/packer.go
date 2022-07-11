package ec2

import (
	"fmt"
	"os"
	"os/exec"
)

type PackerOptions struct {
	BuildPkg string
	AmiName  string
	Arch     string
	OS       string
}

func addArg(args *[]string, arg string) {
	*args = append(*args, "-var")
	*args = append(*args, arg)
}

type AMIArg struct {
	SourceAMIFilter string
	Owner           string
}

// os -> arch -> ami
var packageToAMIArg = map[string]map[string]AMIArg{
	"amzn2":   {"aarch64": AMIArg{SourceAMIFilter: "amzn2-ami-hvm-2.0.*.1-arm64-gp2", Owner: "amazon"}, "x86_64": AMIArg{SourceAMIFilter: "amzn2-ami-hvm-2.0.*.1-x86_64-gp2", Owner: "amazon"}},
	"centos7": {"x86_64": AMIArg{SourceAMIFilter: "CentOS Linux 7 x86_64 HVM EBS *", Owner: "679593333241"}},
}

var osToSSHUsername = map[string]string{
	"amzn2":   "ec2-user",
	"centos7": "centos",
}

var osToPackerfilePrefix = map[string]string{
	"amzn2":   "yum",
	"centos7": "yum",
}

var osToDeviceName = map[string]string{
	"amzn2":   "/dev/xvda",
	"centos7": "/dev/sda1",
}

func CallPacker(opts PackerOptions) error {
	filePath := fmt.Sprintf("packerfiles/aws-%s.pkr.hcl", osToPackerfilePrefix[opts.OS])

	err := exec.Command("packer", "init", filePath).Run()

	if err != nil {
		return err
	}

	args := []string{}

	args = append(args, "build")
	addArg(&args, "build_pkg="+opts.BuildPkg)
	addArg(&args, "ami_name="+opts.AmiName)
	addArg(&args, "arch="+opts.Arch)
	addArg(&args, "source_ami_filter="+packageToAMIArg[opts.OS][opts.Arch].SourceAMIFilter)
	addArg(&args, "ssh_username="+osToSSHUsername[opts.OS])
	addArg(&args, "ami_owner="+packageToAMIArg[opts.OS][opts.Arch].Owner)
	addArg(&args, "device_name="+osToDeviceName[opts.OS])
	args = append(args, filePath)

	cmd := exec.Command("packer", args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()

	return err
}
