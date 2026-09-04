# Use the official lightweight Python image
FROM python:3.14-slim

ARG PORT
ENV PORT=${PORT}

RUN echo "Given port: ${PORT}"

# Set the working directory inside the container
WORKDIR /app

# 1. Install essential system dependencies
# These libraries (libnss3, libatk, etc.) are required by Chromium to run in a headless environment.
# 'curl' is included for health checks or additional downloads if needed.
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    && rm -rf /var/lib/apt/lists/*

# 2. Install 'uv' for Python package and project management
RUN pip install uv

# 3. Handle project dependencies
# Copy only dependency files first to leverage Docker layer caching.
COPY pyproject.toml uv.lock ./

# Synchronize dependencies using the lockfile to ensure a consistent environment.
RUN uv sync --frozen

# 4. Install Playwright browsers and their system-level dependencies.
# We use 'uv run' to ensure the installation happens within the project's virtual environment.
# '--with-deps' ensures that all missing Linux shared libraries for Chromium are present.
RUN uv run playwright install --with-deps chromium

# 5. Copy the rest of the application source code
COPY . .

# Prevents Python from writing .pyc files to disk (keeps the container clean)
ENV PYTHONDONTWRITEBYTECODE=1

# Forces the stdout and stderr streams to be unbuffered (immediate logging in Docker)
ENV PYTHONUNBUFFERED=1

# 6. Execution Command
CMD ["sh", "-c", "uv run gunicorn wsgi:app --workers 3 --timeout 120 --bind 0.0.0.0:$PORT"]