FROM python:3.14-slim

WORKDIR /app

RUN pip install uv

COPY pyproject.toml uv.lock ./
RUN uv sync --frozen

COPY . .

# Build-time args
ARG SERVER_PORT=6742
ARG RUN_COMMAND="gunicorn app.main:app --workers 1"

# Convert them to runtime env vars
ENV SERVER_PORT=${SERVER_PORT}
ENV RUN_COMMAND=${RUN_COMMAND}

EXPOSE ${SERVER_PORT}

CMD sh -c "uv run $RUN_COMMAND -b 0.0.0.0:$SERVER_PORT"