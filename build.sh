#!/bin/bash

REMOTE=origin
TARGET_BRANCH=main
RELEASE_BRANCH=main

build() {
    commit
    push
    merge
}

commit() {
    read -p "enter commit message (default: 🐞 fix:): " commit_msg
    if [ -z "$commit_msg" ]; then
        commit_msg="🐞 fix:"
    else
        commit_msg="🐞 fix: $commit_msg"
    fi

    git add .
    committed_files=$(git diff --staged --name-only --color=always)
    git commit -m "$commit_msg"
    echo "committed files:"
    # echo -e "\e[32m$committed_files\e[0m"
    if [[ "$(uname)" == "Darwin" ]]; then
        tput setaf 3
        echo "$committed_files"
        tput setaf 7
    else
        echo -e "\e[32m$committed_files\e[0m"
    fi
}

push() {
    current_branch=$(git rev-parse --abbrev-ref HEAD)

    echo "pushing to $REMOTE/$current_branch repository..."
    git pull $REMOTE $current_branch

    last_push_commit=$(git rev-parse $REMOTE/$current_branch)
    git push $REMOTE $current_branch
    current_push_commit=$(git rev-parse $REMOTE/$current_branch)

    echo "commits pushed to remote repository:"
    git log --oneline $last_push_commit..$current_push_commit
}

sync() {
    commit
    push
    echo "synced current branch with remote repository."
}

merge() {
    origin_branch=$(git rev-parse --abbrev-ref HEAD)

    git checkout $TARGET_BRANCH
    echo "updating $TARGET_BRANCH branch..."
    git fetch $REMOTE $TARGET_BRANCH
    git reset --hard $REMOTE/$TARGET_BRANCH

    echo "merging $origin_branch into $TARGET_BRANCH..."
    git merge -m "merge $origin_branch into $TARGET_BRANCH" --no-ff $origin_branch

    if git diff --quiet; then
        echo "no conflicts detected."
    else
        echo "conflict detected. rolling back merge..."
        git merge --abort
        echo "merge aborted. conflict detected."
        git status -s | awk '{print $2}' | grep -E "^(U|AA)" | xargs -I {} echo "conflict in file: {}"
        exit 1
    fi

    echo "pushing $TARGET_BRANCH to remote..."
    last_push_commit=$(git rev-parse $REMOTE/$TARGET_BRANCH)
    git push $REMOTE $TARGET_BRANCH
    current_push_commit=$(git rev-parse $REMOTE/$TARGET_BRANCH)
    git log --oneline $last_push_commit..$current_push_commit

    git checkout $origin_branch
}

tag() {
    current_branch=$(git rev-parse --abbrev-ref HEAD)
    if [ "$current_branch" != "$RELEASE_BRANCH" ]; then
        echo "error: you must be on the 'main' branch to create a tag."
        exit 1
    fi

    read -p "enter tag name (format should be v*.*.*): " tag_name
    if [ -z "$tag_name" ] || ! [[ $tag_name =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        echo "please enter a valid tag_name, format should be v*.*.*"
        exit 1
    fi

    echo "updating $RELEASE_BRANCH tags and branch..."
    git fetch --tags
    git fetch $REMOTE $RELEASE_BRANCH
    git reset --hard $REMOTE/$RELEASE_BRANCH

    if git rev-parse -q --verify "refs/tags/$tag_name" >/dev/null; then
        echo "tag $tag_name already exists, please choose a new tag"
        exit 1
    fi
    echo "create $tag_name..."
    git tag -a $tag_name -m "release $tag_name"
    git push origin $tag_name
}

if [ "$1" == "build" ]; then
    build
elif [ "$1" == "commit" ]; then
    commit
elif [ "$1" == "push" ]; then
    push
elif [ "$1" == "sync" ]; then
    sync
elif [ "$1" == "merge" ]; then
    merge
elif [ "$1" == "tag" ]; then
    tag
elif [ "$1" == "help" ]; then
    echo "Usage: $0 {build|commit|push|sync|merge}"
    exit 0
else
    echo "Usage: $0 {build|commit|push|sync|merge}"
    exit 1
fi
