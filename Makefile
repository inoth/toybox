# Makefile

REMOTE := origin
TARGET_BRANCH := main
RELEASE_BRANCH := main

.PHONY: build commit push sync merge tag

build: commit push merge

commit:
	@./build.sh commit

push:
	@./build.sh push

sync:
	@./build.sh sync

merge:
	@./build.sh merge

tag:
	@./build.sh tag

help:
	@echo "Usage: make [target]"
	@echo "Targets:"
	@echo "  build     - Commit and push changes"
	@echo "  commit    - Commit changes"
	@echo "  push      - Push changes to remote repository"
	@echo "  sync      - Sync current branch with remote repository"
	@echo "  merge     - Merge dev branch into main branch"
	@echo "  tag       - Create a new release tag (only on main branch)"
