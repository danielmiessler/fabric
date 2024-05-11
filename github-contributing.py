import sys
import argparse
import subprocess

def get_github_username():
    """Retrieve GitHub username from local Git configuration."""
    result = subprocess.run(['git', 'config', '--get', 'user.name'], capture_output=True, text=True)
    if result.returncode == 0 and result.stdout:
        return result.stdout.strip()
    else:
        raise Exception("Failed to retrieve GitHub username from Git config.")
    
def update_fork():
    # Sync your fork's main branch with the original repository's main branch
    print("Updating fork...")
    subprocess.run(['git', 'fetch', 'upstream'], check=True)  # Fetch the branches and their respective commits from the upstream repository
    subprocess.run(['git', 'checkout', 'main'], check=True)  # Switch to your local main branch
    subprocess.run(['git', 'merge', 'upstream/main'], check=True)  # Merge changes from upstream/main into your local main branch
    subprocess.run(['git', 'push', 'origin', 'main'], check=True)  # Push the updated main branch to your fork on GitHub
    print("Fork updated successfully.")

def create_branch(branch_name):
    print(f"Creating new branch '{branch_name}'...")
    subprocess.run(['git', 'checkout', '-b', branch_name], check=True)
    print(f"Branch '{branch_name}' created and switched to.")

def push_changes(branch_name, commit_message):
    # Push your local changes to your fork on GitHub
    print("Pushing changes to fork...")
    subprocess.run(['git', 'checkout', branch_name], check=True)  # Switch to the branch where your changes are
    subprocess.run(['git', 'add', '.'], check=True)  # Stage all changes for commit
    subprocess.run(['git', 'commit', '-m', commit_message], check=True)  # Commit the staged changes with a custom message
    subprocess.run(['git', 'push', 'fork', branch_name], check=True)  # Push the commit to the same branch in your fork
    print("Changes pushed successfully.")

def create_pull_request(branch_name, pr_title, pr_file):
    # Create a pull request on GitHub using the GitHub CLI
    print("Creating pull request...")
    github_username = get_github_username()
    with open(pr_file, 'r') as file:
        pr_body = file.read()  # Read the PR description from a markdown file
    subprocess.run(['gh', 'pr', 'create',
                    '--base', 'main',
                    '--head', f'{github_username}:{branch_name}',
                    '--title', pr_title,
                    '--body', pr_body], check=True)  # Create a pull request with the specified title and markdown body
    print("Pull request created successfully.")

def main():
    parser = argparse.ArgumentParser(description="Automate your GitHub workflow")
    subparsers = parser.add_subparsers(dest='command', help='Available commands')

    # Subparser for updating fork
    parser_update = subparsers.add_parser('update-fork', help="Update fork with the latest from the original repository")
    
    parser_create_branch = subparsers.add_parser('create-branch', help="Create a new branch")
    parser_create_branch.add_argument('--branch-name', required=True, help="The name for the new branch")

    # Subparser for pushing changes
    parser_push = subparsers.add_parser('push-changes', help="Push local changes to the fork")
    parser_push.add_argument('--branch-name', required=True, help="The name of the branch you are working on")
    parser_push.add_argument('--commit-message', required=True, help="The commit message for your changes")

    # Subparser for creating a pull request
    parser_pr = subparsers.add_parser('create-pr', help="Create a pull request to the original repository")
    parser_pr.add_argument('--branch-name', required=True, help="The name of the branch the pull request is from")
    parser_pr.add_argument('--pr-title', required=True, help="The title of your pull request")
    parser_pr.add_argument('--pr-file', required=True, help="The markdown file path for your pull request description")

    args = parser.parse_args()

    if args.command == 'update-fork':
        update_fork()
    elif args.command == 'create-branch':
        create_branch(args.branch_name)
    elif args.command == 'push-changes':
        push_changes(args.branch_name, args.commit_message)
    elif args.command == 'create-pr':
        create_pull_request(args.branch_name, args.pr_title, args.pr_file)

if __name__ == '__main__':
    main()
