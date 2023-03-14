source provider.env 

rm ~/.terraform.d/plugins/${host_name}/${namespace}/${type}/${version}/${target}/terraform-provider-mattr
cp terraform-provider-mattr ~/.terraform.d/plugins/${host_name}/${namespace}/${type}/${version}/${target}
