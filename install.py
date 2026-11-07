import os
import sys
import subprocess
import tempfile
from pathlib import Path

GIT_REPO = "https://github.com/Asterlenist/kk.git"

def check_root():
    if os.geteuid() != 0:
        print("Error! Run this installer via sudo!", file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    check_root()
    if subprocess.run(["which", "go"], capture_output=True).returncode == 0:
        print("Go is already installed. Let's continue...")
    else:
        print("Go is not installed...", file=sys.stderr)
        sys.exit(1)
    
    source_dir = Path.cwd()
    temporary_directory = None
    if (source_dir / "main.go").is_file():
        print("main.go file found. Don't copy the Git repository...")
    else:
        print("Copying Git repository...")
        temporary_directory = tempfile.TemporaryDirectory(prefix="kk-install-")
        source_dir = Path(temporary_directory.name) / "kk"
        subprocess.run(["git", "clone", GIT_REPO, str(source_dir)], check=True)

    binary_path = source_dir / "kk"
    subprocess.run(["go", "build", "-o", str(binary_path), "."], cwd=source_dir, check=True)
    subprocess.run(["mv", str(binary_path), "/usr/bin/"], check=True)

    print("Installation complete! You can now run 'kk' from anywhere.")