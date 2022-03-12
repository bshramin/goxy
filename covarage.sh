coverage_required=0.3
coverage=$(printf "%.2f\n" $(go test -mod=vendor ./... --cover | awk '{if ($1 != "?") print $5; else print "0.0";}' | sed 's/\%//g' | awk '{s+=$1; n+=1} END {printf "%.5f\n", s/n}'))

printf "Coverage of the whole project: %.2f percent of statements\n" $coverage
if awk "BEGIN {exit !($coverage < $coverage_required)}"; then
  printf "This code does not meet the coverage requirement of %.2f" $coverage_required
  exit 1
fi
