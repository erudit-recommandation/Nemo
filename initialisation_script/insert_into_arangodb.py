# -*- coding: utf-8 -*-

import pandas as pd
from pathlib import Path
from pyArango.connection import *
import numpy as np

directory = Path('./data/erudit')
arangoURL = 'http://localhost:8529'
username = 'root'
password = 'rootpassword'
databaseName = "erudit"
collectionName = "articles"


def convert_df_to_arango(row):
    doc = collection.createDocument()
    for c in cols:
        if row[c] == None:
            pass
        elif c == "sstitrerev":
            sstitrerev = row[c]
            sstitrerev = sstitrerev.split(" •")
            sstitrerev = [rev if rev[0] != " " else rev[1:]
                          for rev in sstitrerev]
            doc[c] = sstitrerev

        else:
            doc[c] = row[c]
    return doc


def initialise_arango_db():
    conn = Connection(arangoURL=arangoURL,
                      username=username,
                      password=password)
    try:
        conn.createDatabase(name=databaseName)
    except:
        print("the database already exist")

    db = conn[databaseName]

    collection = None
    if db.hasCollection(collectionName):
        collection = db[collectionName]
    else:
        collection = db.createCollection(name=collectionName)
        collection.ensureFulltextIndex(fields=["text"])
        collection.ensurePersistentIndex(fields=["idproprio"], unique=True)
        collection.ensurePersistentIndex(fields=["title"])
        collection.ensurePersistentIndex(fields=["author"])

    chunksize = 100
    notice_freq = 1000
    nDocumentAdded = 0

    df = pd.read_csv(path, encoding='utf-8', chunksize=chunksize,
                     sep=';', lineterminator='\n', usecols=cols)

    with df as reader:
        for chunk in reader:
            chunk.reset_index(inplace=True)
            chunk.replace({np.nan: None}, inplace=True)
            for (i, row) in chunk.iterrows():
                doc = convert_df_to_arango(row)

                try:
                    doc.save(waitForSync=False)
                except:
                    print("---an error happen when inserting {}---".format(row))
                finally:
                    nDocumentAdded += 1
                if nDocumentAdded % notice_freq == 0:
                    print("----{} document added----".format(nDocumentAdded))
    print("---- Done -----")
