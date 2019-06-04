package run

import (
	"fmt"
	"github.com/luxas/ignite/pkg/constants"
	"github.com/luxas/ignite/pkg/metadata"
	"github.com/luxas/ignite/pkg/metadata/kernmd"
	"github.com/luxas/ignite/pkg/util"
)

type ImportKernelOptions struct {
	Source      string
	Name        string
	KernelNames []*metadata.Name
}

func ImportKernel(ao *ImportKernelOptions) error {
	if !util.FileExists(ao.Source) {
		return fmt.Errorf("not a kernel image: %s", ao.Source)
	}

	// Create a new ID for the kernel
	kernelID, err := util.NewID(constants.KERNEL_DIR)
	if err != nil {
		return err
	}

	// Verify the name
	name, err := metadata.NewName(ao.Name, &ao.KernelNames)
	if err != nil {
		return err
	}

	md := kernmd.NewKernelMetadata(kernelID, name)

	// Save the metadata
	if err := md.Save(); err != nil {
		return err
	}

	// Perform the copy
	if err := md.ImportKernel(ao.Source); err != nil {
		return err
	}

	fmt.Println(md.ID)

	return nil
}