# Simpl - CLI

# Table of contents
- [Simpl - CLI](#simpl---cli)
- [Notes](#notes)
  - [Autenticazione](#autenticazione)

```bash
simpl-cli --help
simpl-cli login --help
simpl-cli token --help
simpl-cli profile --help
```

# Notes

CLI per l'utilizzo dei servizi di SIMPL da linea di comando.
Una parte necessaria all'utilizzo dei servizi di SIMPL e dover recuperare un token che rappresenta l'utente.

L'utente esegue l'autenticazione usando il protocollo AUTH2 utilizzando il browser per l'autenticazione.

Siccome la CLI è un tool di dev può essere utile permettere l'utilizzo della CLI simulando il comportamento di diversi utenti.

Un utente è caratterizzato da un identitificativo e il server di autenticazione associato.

Ho bisogno delle seguenti funzionalità:

- Autenticazione
- Tokenizzazione

## Autenticazione

Possibilità di definire un utente e registrarsi all'authentication manager.

Mi immagino una funzionalità del tipo:

```bash
simpl-cli auth --user username --authserver https://authserver.com/xxxxx bash
```

La cli permettere all'utente di effettaure l'autenticazione tramite Browser

Se l'autenticazione andrà a buon fine verranno salvate le informazioni di autenticazione che serviranno per la tokenizzazione in altre funzionalità della CLI.
