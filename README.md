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
  awstags [OPTIONS] <delete | get | set>

Help Options:
  -h, --help  Show this help message

Available commands:
  delete  delete all tags of an instance
  get     get tags of an instance, and print to stdout to json
  set     set tags from file/stdin json
```

#### Set

```
% awstags set --help
Usage:
  awstags [OPTIONS] set [set-OPTIONS]

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
% awstags get --help 
Usage:
  awstags [OPTIONS] get [get-OPTIONS]

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
% awstags delete --help
Usage:
  awstags [OPTIONS] delete [delete-OPTIONS]

Help Options:
  -h, --help              Show this help message

[delete command options]
      -p, --profile-name= login using a specific shared credentials profile name
      -k, --key-id=       login using a specific keyId
      -s, --secret-key=   login using a specific secretKey
      -r, --region=       use a specific AWS region
      -i, --instance-id=  required: instance-id to query/set
```

## Use Examples

```bash
# get tags to file
./awstags get -i i-0864c0fb1716d91ca -r us-west-2 > my.json
# delete all tags
./awstags delete -i i-0864c0fb1716d91ca -r us-west-2
# set tags again
./awstags set -i i-0864c0fb1716d91ca -r us-west-2 -f my.json
# set tags - alternative method
cat my.json |./awstags set -i i-0864c0fb1716d91ca -r us-west-2
```
