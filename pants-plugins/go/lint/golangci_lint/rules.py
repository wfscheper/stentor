import textwrap
from dataclasses import dataclass

from go.lint.golangci_lint.skip_field import SkipGolangciLintField
from go.lint.golangci_lint.subsystem import GolangciLint
from pants.backend.go.subsystems import golang
from pants.backend.go.subsystems.golang import GoRoot
from pants.backend.go.target_types import GoPackageSourcesField
from pants.backend.go.util_rules.go_mod import (
    GoModInfo,
    GoModInfoRequest,
    OwningGoMod,
    OwningGoModRequest,
)
from pants.core.goals.lint import LintResult, LintResults, LintTargetsRequest
from pants.core.util_rules.config_files import ConfigFiles, ConfigFilesRequest
from pants.core.util_rules.external_tool import (
    DownloadedExternalTool,
    ExternalToolRequest,
)
from pants.core.util_rules.source_files import SourceFiles, SourceFilesRequest
from pants.core.util_rules.system_binaries import BashBinary
from pants.engine.fs import CreateDigest, Digest, FileContent, MergeDigests
from pants.engine.internals.selectors import Get, MultiGet
from pants.engine.platform import Platform
from pants.engine.process import FallibleProcessResult, Process
from pants.engine.rules import collect_rules, rule
from pants.engine.target import FieldSet, Target
from pants.engine.unions import UnionRule
from pants.util.logging import LogLevel


@dataclass(frozen=True)
class GolangciLintFieldSet(FieldSet):
    required_fields = (GoPackageSourcesField,)

    sources: GoPackageSourcesField

    @classmethod
    def opt_out(cls, tgt: Target) -> bool:
        return tgt.get(SkipGolangciLintField).value


class GolangciLintRequest(LintTargetsRequest):
    field_set_type = GolangciLintFieldSet
    name = GolangciLint.options_scope


@rule(desc="Lint with golangci-lint", level=LogLevel.DEBUG)
async def run_golangci_lint(
    request: GolangciLintRequest,
    golangci_lint: GolangciLint,
    goroot: GoRoot,
    bash: BashBinary,
) -> LintResults:
    if golangci_lint.skip:
        return LintResults([], linter_name=request.name)

    downloaded_golangci_lint, config_files = await MultiGet(
        Get(
            DownloadedExternalTool,
            ExternalToolRequest,
            golangci_lint.get_request(Platform.current),
        ),
        Get(ConfigFiles, ConfigFilesRequest, golangci_lint.config_request()),
    )

    source_files = await Get(
        SourceFiles,
        SourceFilesRequest(field_set.sources for field_set in request.field_sets),
    )

    owning_go_mods = await MultiGet(
        Get(OwningGoMod, OwningGoModRequest(field_set.address))
        for field_set in request.field_sets
    )

    owning_go_mod_addresses = {x.address for x in owning_go_mods}

    go_mod_infos = await MultiGet(
        Get(GoModInfo, GoModInfoRequest(address)) for address in owning_go_mod_addresses
    )

    # golangci-lint requires a absolute path to a cache
    golangci_lint_run_script = FileContent(
        "__run_golangci_lint.sh",
        textwrap.dedent(
            f"""\
            export GOROOT={goroot.path}
            export PATH="${{GOROOT}}/bin"
            sandbox_root="$(/bin/pwd)"
            export GOPATH="${{sandbox_root}})/gopath"
            export GOCACHE="${{sandbox_root}}/gocache"
            export GOLANGCI_LINT_CACHE="$GOCACHE"
            /bin/mkdir -p "$GOPATH" "$GOCACHE"
            exec "$@"
            """
        ).encode("utf-8"),
    )

    golangci_lint_run_script_digest = await Get(
        Digest, CreateDigest([golangci_lint_run_script])
    )

    input_digest = await Get(
        Digest,
        MergeDigests(
            [
                golangci_lint_run_script_digest,
                downloaded_golangci_lint.digest,
                config_files.snapshot.digest,
                source_files.snapshot.digest,
                *(info.digest for info in set(go_mod_infos)),
            ]
        ),
    )

    argv = [
        bash.path,
        golangci_lint_run_script.path,
        downloaded_golangci_lint.exe,
        "run",
    ]
    if golangci_lint.config:
        argv.append(f"--config={golangci_lint.config}")
    elif config_files.snapshot.files:
        argv.append(f"--config={config_files.snapshot.files[0]}")
    else:
        argv.append("--no-config")
    argv.extend(golangci_lint.args)

    process_result = await Get(
        FallibleProcessResult,
        Process(
            argv=argv,
            input_digest=input_digest,
            description="Run `golangci-lint`.",
            level=LogLevel.DEBUG,
        ),
    )

    result = LintResult.from_fallible_process_result(process_result)
    return LintResults([result], linter_name=request.name)


def rules():
    return [
        *collect_rules(),
        *golang.rules(),
        UnionRule(LintTargetsRequest, GolangciLintRequest),
    ]
