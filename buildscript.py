import argparse, subprocess, sys


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "-v",
        "--verbosity",
        type=bool,
        required=False,
        default=False,
        action=argparse.BooleanOptionalAction,
        help="Make test verbose",
    )
    parser.add_argument(
        "-p",
        "--package",
        type=str,
        required=False,
        default="./...",
        help="Package name",
    )
    parser.add_argument(
        "-t",
        "--target",
        type=str,
        required=False,
        default=None,
        help="Target test name",
    )
    parser.add_argument(
        "--showNoTestFilesWarning",
        type=bool,
        required=False,
        default=False,
        action=argparse.BooleanOptionalAction,
        help="Shows no test files warning",
    )
    args = parser.parse_args()

    command = "set -o pipefail; go test -race"
    command = command + " " + args.package
    if args.verbosity:
        command = command + " -v"
    if args.target:
        command = command + " -run " + args.target
    if not args.showNoTestFilesWarning:
        command = command + ' | { grep -v "no test files"; true; }'
    process = subprocess.run(
        command,
        shell = True,
        executable = "/bin/bash",
        check = False
    )
    sys.exit(process.returncode)


if __name__ == "__main__":
    main()