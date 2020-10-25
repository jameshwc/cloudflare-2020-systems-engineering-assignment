BINARY="simple-http"
DEFAULT_URL="https://cloudflare-2020-general-engineering-assignment.jameshwc.workers.dev/links"
PROFILE_NUM="100"
STRESS_NUM="100"
STRESS_CONCURRENCY="100"

.PHONY: default single profile stress clean
default:
	go build -o ${BINARY}

single:
	./${BINARY} -url ${DEFAULT_URL}

profile:
	./${BINARY} -url ${DEFAULT_URL} -profile ${PROFILE_NUM}

stress:
	./${BINARY} -url ${DEFAULT_URL} -profile ${STRESS_NUM} -c ${STRESS_CONCURRENCY}

clean:
	rm -f ${BINARY}