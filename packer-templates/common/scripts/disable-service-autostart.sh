#!/bin/bash
export DEBIAN_FRONTEND=noninteractive
set -euxo pipefail

echo -e '#!/bin/sh\nexit 101' > /usr/sbin/policy-rc.d
chmod +x /usr/sbin/policy-rc.d