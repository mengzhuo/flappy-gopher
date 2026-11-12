from pathlib import Path
import yaml

ROOT = Path(__file__).resolve().parents[1]
WORKFLOWS = ROOT / ".github" / "workflows"
GO_COMMANDS = {"go test -count=1 ./...", "go vet ./..."}


def load_workflow(name):
    with (WORKFLOWS / name).open(encoding="utf-8") as stream:
        return yaml.safe_load(stream)


def run_steps(job):
    return {step.get("run") for step in job["steps"] if "run" in step}


def validate_test(workflow):
    assert workflow["name"] == "Test"
    assert workflow["on"]["push"]["branches"] == ["main"]
    assert workflow["on"]["pull_request"]["branches"] == ["main"]
    assert workflow["permissions"]["contents"] == "read"
    job = workflow["jobs"]["test"]
    assert job["strategy"]["fail-fast"] is False
    assert set(job["strategy"]["matrix"]["os"]) == {
        "windows-latest", "macos-latest", "ubuntu-latest"
    }
    assert job["runs-on"] == "${{ matrix.os }}"
    assert GO_COMMANDS <= run_steps(job)
    uses = [step["uses"] for step in job["steps"] if "uses" in step]
    assert uses == ["actions/checkout@v7", "actions/setup-go@v7"]


if __name__ == "__main__":
    validate_test(load_workflow("test.yml"))
    print("test.yml: OK")
