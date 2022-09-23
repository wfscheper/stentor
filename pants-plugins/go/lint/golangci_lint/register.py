from go.lint.golangci_lint import skip_field
from go.lint.golangci_lint.rules import rules as golangci_lint_rules
from pants.backend.experimental.go.register import rules as go_rules


def rules():
    return [
        *go_rules(),
        *golangci_lint_rules(),
        *skip_field.rules(),
    ]
