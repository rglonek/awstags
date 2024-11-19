# AWS-Tags

Simple application for getting, setting and deleting tags on EC2 instances.

## Install

1. Head to the [releases](https://github.com/rglonek/awstags/releases) page and download the binary you want to use.
2. Rename it to `awstags` for ease of typing. Ex: `mv awstags.darwin.amd64 awstags`
3. Make it executable: `chmod 755 awstags`
4. Run it

## Usage

### Help pages

#### Main

```
% awstags --help
Usage:
  awstags [OPTIONS] <ec2 | efs>

Help Options:
  -h, --help  Show this help message

Available commands:
  ec2  operate on ec2 tags
  efs  operate on efs tags
```

```
% awstags ec2 --help
Usage:
  awstags [OPTIONS] ec2 <delete | get | list | set>

Help Options:
  -h, --help  Show this help message

Available commands:
  delete  delete all tags of an instance
  get     get tags of an instance, and print to stdout to json
  list    list all ec2 instances
  set     set tags from file/stdin json
```

#### Set

```
% awstags ec2 set --help
Usage:
  awstags [OPTIONS] ec2 set [set-OPTIONS]

Help Options:
  -h, --help              Show this help message

[set command options]
      -p, --profile-name= login using a specific shared credentials profile name
      -k, --key-id=       login using a specific keyId
      -s, --secret-key=   login using a specific secretKey
      -r, --region=       use a specific AWS region
      -i, --instance-id=  required: instance-id to query/set
      -f, --filename=     filename of the file to read tags from; will read os.Stdin if not specified
```

#### Get

```
% awstags ec2 get --help 
Usage:
  awstags [OPTIONS] ec2 get [get-OPTIONS]

Help Options:
  -h, --help              Show this help message

[get command options]
      -p, --profile-name= login using a specific shared credentials profile name
      -k, --key-id=       login using a specific keyId
      -s, --secret-key=   login using a specific secretKey
      -r, --region=       use a specific AWS region
      -i, --instance-id=  required: instance-id to query/set
```

#### Delete

```
% awstags ec2 delete --help
Usage:
  awstags [OPTIONS] ec2 delete [delete-OPTIONS]

Help Options:
  -h, --help              Show this help message

[delete command options]
      -p, --profile-name= login using a specific shared credentials profile name
      -k, --key-id=       login using a specific keyId
      -s, --secret-key=   login using a specific secretKey
      -r, --region=       use a specific AWS region
      -i, --instance-id=  required: instance-id to query/set
```

#### List

```
% awstags ec2 list --help
Usage:
  awstags [OPTIONS] ec2 list [list-OPTIONS]

Help Options:
  -h, --help              Show this help message

[list command options]
      -p, --profile-name= login using a specific shared credentials profile name
      -k, --key-id=       login using a specific keyId
      -s, --secret-key=   login using a specific secretKey
      -r, --region=       use a specific AWS region
      -t, --with-tags     list instances with tags - produces long json
```

## Use Examples

```bash
# get tags to file
./awstags ec2 get -i i-0864c0fb1716d91ca -r us-west-2 > my.json
# delete all tags
./awstags ec2 delete -i i-0864c0fb1716d91ca -r us-west-2
# set tags again
./awstags ec2 set -i i-0864c0fb1716d91ca -r us-west-2 -f my.json
# set tags - alternative method
cat my.json |./awstags ec2 set -i i-0864c0fb1716d91ca -r us-west-2
# list all efs volumes
./awstags efs list -r us-west-2
# list all efs volumes with tags as json
./awstags efs list -r us-west-2 -t

# get all efs volumes and tags in separate files - the inefficient method
./awstags efs list -r us-west-2 |while read efsid; do
  ./awstags efs get -r us-west-2 -e $efsid > $efsid.json
  sleep 1 # sleep implemented to avoid hitting AWS API limits
done

# get all efs volumes and tags in separate files - the good method
./awstags efs list -r us-west-2 -t |jq -c 'to_entries[]' | while read -r entry; do
  key=$(echo "$entry" | jq -r '.key')
  value=$(echo "$entry" | jq '.value')
  echo "$value" > "$key.json"
done
```
