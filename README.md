# CBP60-DB Web Application Source Code
## A Plant Kingdom-Wide Structural Database of the CBP60 Protein Family.

This database features more than 2,000 proteins of the CBP60 family in the plant kingdom. Deep learning-predicted structures, cDNA/amino acid sequences, species of origin, and access to external protein motif/domain analyses are provided.

If you find the contents of our database or the software useful please cite our [biorxiv preprint](https://www.biorxiv.org/content/10.1101/2022.07.07.499200v1).

### Important Links:
- [Link to live web application](https://cbp60db.wlu.ca/)
- [Link to preprint](https://www.biorxiv.org/content/10.1101/2022.07.07.499200v1)
- [Link to Google Colab notebook for protein clustering](https://colab.research.google.com/drive/1LOZY33CSO5-PdJAdDApyPlfxUu4DHjcW)

### Local Setup
Requires golang to be installed [https://go.dev/](https://go.dev/)
```sh
# clone the repository
git clone https://github.com/KeaunAmani/cbp60db.git

# go into main directory
cd cbp60db

# extract the TM-Align matrix
tar -xf tm-matrix.json.tar.bz2

# build binary
go build main.go

# execute binary locally for access via port 8002
./main
```

If you choose to daemonize the service with systemd use the `cbb60db.service` file and configure each line that contains a `#TODO`. Descriptions are provided for each line.