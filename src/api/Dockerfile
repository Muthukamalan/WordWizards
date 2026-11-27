FROM ghcr.io/astral-sh/uv:python3.12-bookworm-slim

WORKDIR /app

ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1
ENV PATH="/app/.venv/bin:$PATH"

COPY src/requirements.txt .

# RUN pip install --no-cache-dir -r requirements.txt
RUN uv venv /app/.venv && uv pip install --no-cache -r requirements.txt

COPY src/ .

RUN useradd --create-home appuser && chown -R appuser:appuser /app

USER appuser

EXPOSE 8080

CMD ["gunicorn", "--bind", "0.0.0.0:8080", "--workers", "2", "--threads", "4", "--timeout", "30", "app:app"]
