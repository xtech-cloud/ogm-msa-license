package config

const defaultYAML string = `
service:
    name: xtc.ogm.license
    address: :18804
    ttl: 15
    interval: 10
logger:
    level: trace
    dir: /var/log/msa/
database:
    driver: sqlite
    mysql:
        address: 127.0.0.1:3306
        user: root
        password: mysql@XTC
        db: ogm
    sqlite:
        path: /tmp/msa-license.db
`
