#!/usr/bin/env bash
# Print per-language coverage summary and exit non-zero if any language < 80%.
set -euo pipefail

THRESHOLD=80
REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"

LANGUAGES=(python dart go cpp)
LCOV_PATHS=(
  "python/coverage/lcov.info"
  "dart/coverage/lcov.info"
  "go/coverage/lcov.info"
  "cpp/coverage/lcov.info"
)

# Column widths
COL_LANG=10
COL_LINES=8
COL_HIT=8
COL_COV=10
COL_STATUS=10

divider() {
  printf '%*s\n' $((COL_LANG + COL_LINES + COL_HIT + COL_COV + COL_STATUS + 4)) '' | tr ' ' '─'
}

header() {
  printf "%-${COL_LANG}s  %${COL_LINES}s  %${COL_HIT}s  %${COL_COV}s  %s\n" \
    "Language" "Lines" "Hit" "Coverage" "Status"
  divider
}

total_lf=0
total_lh=0
any_fail=0

header

for i in "${!LANGUAGES[@]}"; do
  lang="${LANGUAGES[$i]}"
  lcov_path="$REPO_ROOT/${LCOV_PATHS[$i]}"

  if [[ ! -f "$lcov_path" ]]; then
    printf "%-${COL_LANG}s  %${COL_LINES}s  %${COL_HIT}s  %${COL_COV}s  %s\n" \
      "$lang" "-" "-" "-" "MISSING"
    any_fail=1
    continue
  fi

  # Sum LF/LH records directly from lcov.info
  lf=$(grep -c '^DA:' "$lcov_path" || true)
  lh=$(awk -F: '/^DA:/{split($2,a,","); if(a[2]>0) c++} END{print c+0}' "$lcov_path")

  # Re-sum proper LF/LH from summary records (more accurate)
  lf=$(awk -F: '/^LF:/{s+=$2} END{print s+0}' "$lcov_path")
  lh=$(awk -F: '/^LH:/{s+=$2} END{print s+0}' "$lcov_path")

  total_lf=$((total_lf + lf))
  total_lh=$((total_lh + lh))

  if [[ "$lf" -eq 0 ]]; then
    pct="0.0"
  else
    pct=$(awk "BEGIN{printf \"%.1f\", ($lh/$lf)*100}")
  fi

  pct_int=$(awk "BEGIN{printf \"%d\", int($pct)}")
  pct_display="${pct}%"

  if [[ "$pct_int" -ge "$THRESHOLD" ]]; then
    status="PASS"
  else
    status="FAIL"
    any_fail=1
  fi

  printf "%-${COL_LANG}s  %${COL_LINES}s  %${COL_HIT}s  %${COL_COV}s  %s\n" \
    "$lang" "$lf" "$lh" "$pct_display" "$status"
done

divider

if [[ "$total_lf" -eq 0 ]]; then
  combined_pct="0.0"
else
  combined_pct=$(awk "BEGIN{printf \"%.1f\", ($total_lh/$total_lf)*100}")
fi

printf "%-${COL_LANG}s  %${COL_LINES}s  %${COL_HIT}s  %${COL_COV}s\n" \
  "combined" "$total_lf" "$total_lh" "${combined_pct}%"

if [[ "$any_fail" -ne 0 ]]; then
  echo ""
  echo "FAIL: one or more languages are below the ${THRESHOLD}% coverage threshold."
  exit 1
fi
