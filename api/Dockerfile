FROM python:3.14-slim

WORKDIR /app

# Instalar uv
RUN pip install uv

# Copiar dependencias primero (cache optimizado)
COPY pyproject.toml uv.lock ./

# Instalar dependencias con uv
RUN uv sync --frozen

# Copiar el resto del código
COPY . .

# Variables de entorno
ENV SERVER_PORT=6742
ENV API_SECURITY=true

CMD ["sh", "-c", "uv run gunicorn wsgi:app --workers 1 --bind 0.0.0.0:$SERVER_PORT"]