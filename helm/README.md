# ARYZONA - HELM CHART

Deploy Aryzona with from a helm chart.

# Usage

First, add the helm repo to your helm client:

> `helm repo add dbcafe https://code.db.cafe/api/packages/pauloo27/helm` 

> `helm repo update`

Then, generate the default values.yml: 

> `helm show values aryona/aryzona > my-values-file.yaml`

Now, you can either create a secret with the bot config file or fill the config inside the
values file (under the `config` section).

If you want to use the secret, create it with the following command:

> `kubectl create secret generic <secret-name> --from-file ./config.yml -n <namespace>`

After having all the values in the values file as you wish, install the chart:

> `helm install -f my-values-file.yaml aryona aryona/aryzona -n <namespace>`
