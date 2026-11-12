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


def validate_release(workflow):
    assert workflow["name"] == "Release"
    assert workflow["on"]["push"]["tags"] == ["v*"]
    build = workflow["jobs"]["build"]
    assert build["permissions"]["contents"] == "read"
    assert build["strategy"]["fail-fast"] is False
    matrix = build["strategy"]["matrix"]["include"]
    assert {item["platform"] for item in matrix} == {
        "windows/amd64", "windows/arm64", "darwin/universal",
        "linux/amd64", "linux/arm64", "linux/riscv64",
    }
    assert {item["kind"] for item in matrix} == {"windows", "macos", "linux"}
    for item in matrix:
        assert item["output"] == item["platform"].replace("/", "-")
    assert GO_COMMANDS <= run_steps(build)
    uploads = [
        step for step in build["steps"]
        if step.get("uses") == "actions/upload-artifact@v7"
    ]
    assert len(uploads) == 1
    assert uploads[0]["with"]["if-no-files-found"] == "error"
    assert uploads[0]["with"]["name"] == "release-${{ matrix.output }}"
    publish = workflow["jobs"]["publish"]
    assert publish["needs"] == "build"
    assert publish["permissions"]["contents"] == "write"
    verify = next(
        step["run"] for step in publish["steps"]
        if "sha256sum" in step.get("run", "")
    )
    assert "windows=(package/*windows-*.zip)" in verify
    assert 'test "${#windows[@]}" -eq 2' in verify
    assert 'test "${#macos[@]}" -eq 1' in verify
    assert 'test "${#linux_tar[@]}" -eq 3' in verify
    assert 'test "${#linux_deb[@]}" -eq 3' in verify
    uses = [step["uses"] for step in publish["steps"] if "uses" in step]
    assert uses == ["actions/download-artifact@v8", "softprops/action-gh-release@v3"]
    release = publish["steps"][-1]["with"]
    assert release["files"] == "package/*"
    assert release["generate_release_notes"] is True


if __name__ == "__main__":
    validate_test(load_workflow("test.yml"))
    print("test.yml: OK")
    validate_release(load_workflow("release.yml"))
    print("release.yml: OK")
