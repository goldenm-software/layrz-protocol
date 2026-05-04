#!/usr/bin/env bash
# Convert Go coverage profile (coverage.out) to lcov.info format.
# Usage: gocov2lcov.sh coverage/coverage.out > coverage/lcov.info
set -euo pipefail

PROFILE="${1:-}"
if [[ -z "$PROFILE" ]]; then
  echo "Usage: $0 <coverage.out>" >&2
  exit 1
fi

# Skip the "mode:" header line, then group by source file.
awk '
/^mode:/ { next }
{
  # line format: file.go:startline.col,endline.col stmts count
  split($1, loc, ":")
  file = loc[1]
  split(loc[2], range, ",")
  startline = int(range[1])
  count = int($3)

  files[file] = 1
  if (!(file SUBSEP startline in hits)) {
    hits[file SUBSEP startline] = 0
    order[file][length(order[file])+1] = startline
  }
  # A line covered by any block with count>0 is "hit"
  if (count > 0) hits[file SUBSEP startline] += count
}
END {
  for (file in files) {
    print "SF:" file
    lf = 0; lh = 0
    for (i = 1; i <= length(order[file]); i++) {
      ln = order[file][i]
      cnt = hits[file SUBSEP ln]
      print "DA:" ln "," cnt
      lf++
      if (cnt > 0) lh++
    }
    print "LF:" lf
    print "LH:" lh
    print "end_of_record"
  }
}
' "$PROFILE"
