# gist that I found in "https://stackoverflow.com/questions/5014632/how-can-i-parse-a-yaml-file-from-a-linux-shell-script"
function parse_yaml {
   local prefix=$2
   local s='[[:space:]]*' w='[a-zA-Z0-9_]*' fs=$(echo @|tr @ '\034')
   sed -ne "s|^\($s\):|\1|" \
        -e "s|^\($s\)\($w\)$s:$s[\"']\(.*\)[\"']$s\$|\1$fs\2$fs\3|p" \
        -e "s|^\($s\)\($w\)$s:$s\(.*\)$s\$|\1$fs\2$fs\3|p"  $1 |
   awk -F$fs '{
      indent = length($1)/2;
      vname[indent] = $2;
      for (i in vname) {if (i > indent) {delete vname[i]}}
      if (length($3) > 0) {
         vn=""; for (i=0; i<indent; i++) {vn=(vn)(vname[i])("_")}
         printf("%s%s%s=\"%s\"\n", "'$prefix'",vn, $2, $3);
      }
   }'
}


# Creation/Assertion of fsdiff pod

podName=$(eval echo `parse_yaml test-network-function/tnf_config.yml | grep fsDiffMasterContainer | grep podName | cut -d "_" -f 3 | cut -d "=" -f 2`)
podNamespace=$(eval echo `parse_yaml test-network-function/tnf_config.yml | grep fsDiffMasterContainer | grep namespace | cut -d "_" -f 3 | cut -d "=" -f 2`)

#$podName=node-master

x="fail"
oc get pod -n $podNamespace $podName &>/dev/null && x="success"
if [ x == "success" ]; then
    echo "fsDiffMasterPod already exists, no reason to recreate"
else
    oc apply -f $TNF_PARTNER_SRC_DIR/local-test-infra/fsdiff-pod.yaml
fi

# Creation/Assertion of partner pod

podName=$(eval echo `parse_yaml test-network-function/tnf_config.yml | grep testOrchestrator | grep podName | cut -d "_" -f 3 | cut -d "=" -f 2`)
podNamespace=$(eval echo `parse_yaml test-network-function/tnf_config.yml | grep testOrchestrator | grep namespace | cut -d "_" -f 3 | cut -d "=" -f 2`)

x="fail"
oc get pod -n $podNamespace $podName &>/dev/null && x="success"
if [ x == "success" ]; then
    echo "partner pod already exists, no reason to recreate"
else
    oc apply -f $TNF_PARTNER_SRC_DIR/local-test-infra/local-partner-pod.yaml
fi
