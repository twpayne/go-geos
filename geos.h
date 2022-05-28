#ifndef GEOS_H
#define GEOS_H

#define GEOS_USE_ONLY_R_API
#include <geos_c.h>

int c_GEOSCoordSeq_getFlatCoords_r(GEOSContextHandle_t handle, const GEOSCoordSequence *s, unsigned int size,
                                   unsigned int dims, double *flatCoords);
void c_GEOSGeomBounds_r(GEOSContextHandle_t handle, const GEOSGeometry *g, double *minX, double *minY, double *maxX,
                        double *maxY);
int c_GEOSGeomGetInfo_r(GEOSContextHandle_t handle, const GEOSGeometry *g, int *typeID, int *numGeometries,
                        int *numPoints, int *numInteriorRings);
void c_errorMessageHandler(const char *message, void *userdata);
GEOSCoordSequence *c_newGEOSCoordSeqFromFlatCoords_r(GEOSContextHandle_t handle, unsigned int size, unsigned int dims,
                                                     const double *flatCoords);
GEOSGeometry *c_newGEOSGeomFromBounds_r(GEOSContextHandle_t handle, int *typeID, double minX, double minY, double maxX,
                                        double maxY);

#if GEOS_VERSION_MAJOR < 3 || (GEOS_VERSION_MAJOR == 3 && GEOS_VERSION_MINOR < 10)
GEOSGeometry *GEOSDensify_r(GEOSContextHandle_t handle, const GEOSGeometry *g, double tolerance);
#endif

#endif