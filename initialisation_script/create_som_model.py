import joblib
from gnomonics.preprocessing import lemmatizing_text
import numpy as np
import pandas as pd
import scipy
import sys
from sompy.sompy import SOMFactory
from pathlib import Path
import os
import matplotlib.pyplot as plt
import json

directory = Path('./data/erudit')


def create_som_model(n_rows=100, number_of_som_cells=300):
    path = os.path.join(directory, 'doc_countvectors', 'doc_countvectors.npz')
    X_arr = scipy.sparse.load_npz(path).toarray()

    # "/Users/michael/Dropbox/Data_Projects/Erudit/Data/doc_countvectors/som_grid_bmus.npy"
    path = os.path.join(directory, 'doc_countvectors', 'som_grid_bmus.npy')
    som_grid_bmus = np.load(path)
    X_arr_bmus = np.append(X_arr, np.array([som_grid_bmus]).T, axis=1)

    df_arr = pd.DataFrame(X_arr_bmus)
    df_arr.head()

    sel_arr = []

    for x in set(X_arr_bmus[:, X_arr_bmus.shape[1]-1]):
        df_sel = df_arr[df_arr.iloc[:, X_arr_bmus.shape[1] - 1] == x][:4]
        for y in df_sel.iterrows():
            sel_arr.append(y[1])
        sys.stdout.write("\rBMU %s was processed." % int(x))
        sys.stdout.flush()

    sel_arr = np.asarray(sel_arr)
    np.random.shuffle(sel_arr)

    mapsize = [1, number_of_som_cells]  # [1,10000]
    vectors = sel_arr
    verbose = 'info'

    print("... learning vector dimensionality (SOM)")
    som = SOMFactory().build(vectors, mapsize=[mapsize[0], mapsize[1]])
    # ,train_rough_len=100,train_finetune_len=200)
    som.train(n_job=1, verbose=verbose)

    topographic_error = som.calculate_topographic_error()
    quantization_error = som.calculate_quantization_error()
    print("Topographic error = %s\n Quantization error = %s" %
          (topographic_error, quantization_error))

    path = os.path.join(directory, 'doc_countvectors', 'som_mask_key')
    joblib.dump(som, path, compress=True)
