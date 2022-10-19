package instruction

import (
	"context"

	"github.com/werf/werf/pkg/build/stage"
	"github.com/werf/werf/pkg/config"
	"github.com/werf/werf/pkg/container_backend"
	backend_instruction "github.com/werf/werf/pkg/container_backend/instruction"
	"github.com/werf/werf/pkg/dockerfile"
	dockerfile_instruction "github.com/werf/werf/pkg/dockerfile/instruction"
	"github.com/werf/werf/pkg/util"
)

type Shell struct {
	*Base[*dockerfile_instruction.Shell]
}

func NewShell(name stage.StageName, i *dockerfile.DockerfileStageInstruction[*dockerfile_instruction.Shell], dependencies []*config.Dependency, hasPrevStage bool, opts *stage.BaseStageOptions) *Shell {
	return &Shell{Base: NewBase(name, i, dependencies, hasPrevStage, opts)}
}

func (stage *Shell) GetDependencies(ctx context.Context, c stage.Conveyor, cb container_backend.ContainerBackend, prevImage, prevBuiltImage *stage.StageImage, buildContextArchive container_backend.BuildContextArchiver) (string, error) {
	var args []string
	args = append(args, stage.instruction.Data.Name())
	args = append(args, stage.instruction.Data.Shell...)
	return util.Sha256Hash(args...), nil
}

func (stage *Shell) PrepareImage(ctx context.Context, c stage.Conveyor, cb container_backend.ContainerBackend, prevBuiltImage, stageImage *stage.StageImage, buildContextArchive container_backend.BuildContextArchiver) error {
	return stage.Base.prepareInstruction(ctx, stageImage, buildContextArchive, backend_instruction.NewShell(*stage.instruction.Data))
}
