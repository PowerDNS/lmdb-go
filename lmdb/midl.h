/* lmdb-go v2 midl.h shim: selects the pristine upstream midl.h matching
 * LMDBGO_STREAM. Only the vendored engine sources include this.
 */
#ifndef LMDBGO_MIDL_H_SHIM
#define LMDBGO_MIDL_H_SHIM

#if !defined(LMDBGO_STREAM)
# error "midl.h is internal to the vendored engine builds (LMDBGO_STREAM not set)"
#elif LMDBGO_STREAM == 9
# include "midl_lmdb09.h"
#elif LMDBGO_STREAM == 10
# include "midl_lmdb10.h"
#else
# error "invalid LMDBGO_STREAM (must be 9 or 10)"
#endif

#endif /* LMDBGO_MIDL_H_SHIM */
