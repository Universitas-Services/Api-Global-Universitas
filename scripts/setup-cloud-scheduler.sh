#!/usr/bin/env bash
# Crea (o actualiza) 2 jobs de Cloud Scheduler que despiertan Cloud Run
# y llaman POST /api/v1/economia/bcv/capturar a las 07:00 y 17:00 America/Caracas.
#
# Requisitos:
#   - gcloud autenticado
#   - Variable BCV_CRON_SECRET definida en Cloud Run (mismo valor que CRON_SECRET abajo)
#   - API Cloud Scheduler habilitada
#
# Uso:
#   export CRON_SECRET="tu-secreto-fuerte"
#   bash scripts/setup-cloud-scheduler.sh

set -euo pipefail

PROJECT_ID="${PROJECT_ID:-agente-manual-contrataciones}"
REGION="${REGION:-us-central1}"
SERVICE_URL="${SERVICE_URL:-https://api-global-universitas-693924722323.us-central1.run.app}"
CRON_SECRET="${CRON_SECRET:?Define CRON_SECRET (mismo valor que BCV_CRON_SECRET en Cloud Run)}"

gcloud config set project "$PROJECT_ID"
gcloud services enable cloudscheduler.googleapis.com --project="$PROJECT_ID"

URI="${SERVICE_URL}/api/v1/economia/bcv/capturar"

create_or_update() {
  local name="$1"
  local schedule="$2"

  if gcloud scheduler jobs describe "$name" --location="$REGION" --project="$PROJECT_ID" >/dev/null 2>&1; then
    echo "Actualizando job $name ..."
    gcloud scheduler jobs update http "$name" \
      --location="$REGION" \
      --project="$PROJECT_ID" \
      --schedule="$schedule" \
      --time-zone="America/Caracas" \
      --uri="$URI" \
      --http-method=POST \
      --headers="X-Cron-Secret=${CRON_SECRET},Content-Type=application/json" \
      --attempt-deadline=120s
  else
    echo "Creando job $name ..."
    gcloud scheduler jobs create http "$name" \
      --location="$REGION" \
      --project="$PROJECT_ID" \
      --schedule="$schedule" \
      --time-zone="America/Caracas" \
      --uri="$URI" \
      --http-method=POST \
      --headers="X-Cron-Secret=${CRON_SECRET},Content-Type=application/json" \
      --attempt-deadline=120s
  fi
}

# Cron: minuto hora * * *  (hora Caracas)
create_or_update "bcv-captura-07" "0 7 * * *"
create_or_update "bcv-captura-17" "0 17 * * *"

echo ""
echo "Listo. Jobs:"
gcloud scheduler jobs list --location="$REGION" --project="$PROJECT_ID"

echo ""
echo "Probar manualmente un job:"
echo "  gcloud scheduler jobs run bcv-captura-07 --location=$REGION --project=$PROJECT_ID"
