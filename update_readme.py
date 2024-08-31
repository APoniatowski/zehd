import json


def parse_benchmark_results(file_path):
    with open(file_path, 'r') as f:
        data = [json.loads(line) for line in f]

    results = []
    for entry in data:
        if entry.get("Action") == "output" and "Benchmark" in entry.get("Test", ""):
            results.append(entry)

    return results


def update_readme(benchmark_results, readme_path):
    table_header = "| Package | Test | Runs | ns/op | B/op | allocs/op |\n| --- | --- | --- | --- | --- | --- |\n"
    table_rows = []

    for result in benchmark_results:
        output = result["Output"].strip()
        parts = output.split()

        if len(parts) < 7:
            continue

        test_name = result["Test"]
        package_name = result["Package"]

        if len(parts) >= 7 and parts[0].startswith("Benchmark"):
            runs = parts[1]
            ns_op = parts[2]
            ns_label = parts[3]
            b_op = parts[4]
            b_label = parts[5]
            allocs_op = parts[6]

            if ns_label == "ns/op" and b_label == "B/op":
                row = f"| {package_name} | {test_name} | {
                    runs} | {ns_op} | {b_op} | {allocs_op} |"
                table_rows.append(row)

    table_content = table_header + "\n".join(table_rows)

    with open(readme_path, 'r') as f:
        readme_content = f.read()

    if "## Benchmark Results" in readme_content:
        readme_content = readme_content.split("## Benchmark Results")[0]

    new_readme_content = readme_content + "\n## Benchmark Results\n" + table_content

    with open(readme_path, 'w') as f:
        f.write(new_readme_content)


def main():
    benchmark_results = parse_benchmark_results('benchmark_results.json')
    update_readme(benchmark_results, 'README.md')


if __name__ == "__main__":
    main()
