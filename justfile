default:
  @just --list

# Recipe to run your development environment commands
dev:
  @echo "Entering Nix development environment 🧰 ..."
  @nix develop --extra-experimental-features nix-command --extra-experimental-features flakes

# Recipe to run pre-commit hooks
precommit:
  @echo "Running pre-commit hooks..."
  pre-commit run --all-files
