FROM python:3.12.2-slim

# Install required packages
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        git \
        build-essential \
        ffmpeg \
    && rm -rf /var/lib/apt/lists/*

# Install pipx
RUN python3 -m pip install --upgrade pip \
    && python3 -m pip install --user pipx \
    && python3 -m pipx ensurepath

# Set up work directory
WORKDIR /app

# Clone the repository and install its dependencies
RUN git clone https://github.com/danielmiessler/fabric.git \
    && python3 -m pipx install ./fabric

# Set the entrypoint
ENTRYPOINT ["/usr/bin/bash"]