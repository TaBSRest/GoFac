import argparse

def adjust_coverage(new_coverage):
    try:
        if not new_coverage.endswith("%"):
            new_coverage += "%"
        
        coverage_value = float(new_coverage.strip("%"))
        if coverage_value < 0 or coverage_value > 100:
            raise ValueError("The code coverage must be between 0% and 100%.")
        
        with open("coverage/LineCoverage.txt", "w") as file:
            file.write(new_coverage)
        print(f"The code coverage has been adjusted to {new_coverage}.")
    
    except ValueError as e:
        print(f"Error: {e}")
        exit(1)

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "coverage",
        type=str,
        help="The new code coverage (e.g., '85.0%' or '85.0')",
    )
    args = parser.parse_args()

    adjust_coverage(args.coverage)

if __name__ == "__main__":
    main()