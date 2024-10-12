#ifndef GEOS_H
#define GEOS_H

#include <stdint.h>

#define GEOS_USE_ONLY_R_API
#include <geos_c.h>

uintptr_t c_GEOSGeom_getUserData_r(GEOSContextHandle_t handle,
                                   const GEOSGeometry *g);
void c_GEOSGeom_setUserData_r(GEOSContextHandle_t handle, GEOSGeometry *g,
                              uintptr_t userdata);
void c_GEOSGeomBounds_r(GEOSContextHandle_t handle, const GEOSGeometry *g,
                        double *minX, double *minY, double *maxX, double *maxY);
int c_GEOSGeomGetInfo_r(GEOSContextHandle_t handle, const GEOSGeometry *g,
                        int *typeID, int *numGeometries, int *numPoints,
                        int *numInteriorRings);
void c_errorMessageHandler(const char *message, void *userdata);
GEOSCoordSequence *c_newGEOSCoordSeqFromFlatCoords_r(GEOSContextHandle_t handle,
                                                     unsigned int size,
                                                     unsigned int dims,
                                                     const double *flatCoords);
GEOSGeometry *c_newGEOSGeomFromBounds_r(GEOSContextHandle_t handle, int *typeID,
                                        double minX, double minY, double maxX,
                                        double maxY);
int c_GEOSSTRtree_distance_callback(const void *item1, const void *item2,
                                    double *distance, void *userdata);
void c_GEOSSTRtree_query_callback(void *elem, void *userdata);
GEOSGeometry *c_GEOSMakeValidWithParams_r(GEOSContextHandle_t handle,
                                          const GEOSGeometry *g,
                                          enum GEOSMakeValidMethods method,
                                          int keepCollapsed);

#if GEOS_VERSION_MAJOR < 3 ||                                                  \
    (GEOS_VERSION_MAJOR == 3 && GEOS_VERSION_MINOR < 11)

GEOSGeometry *GEOSConcaveHull_r(GEOSContextHandle_t handle,
                                const GEOSGeometry *g, double ratio,
                                unsigned int allowHoles);

#endif

#if GEOS_VERSION_MAJOR < 3 ||                                                  \
    (GEOS_VERSION_MAJOR == 3 && GEOS_VERSION_MINOR < 12)

GEOSGeometry *GEOSConcaveHullByLength_r(GEOSContextHandle_t handle,
                                        const GEOSGeometry *g, double ratio,
                                        unsigned int allowHoles);

char GEOSPreparedContainsXY_r(GEOSContextHandle_t handle,
                              const GEOSPreparedGeometry *pg1, double x,
                              double y);

char GEOSPreparedIntersectsXY_r(GEOSContextHandle_t handle,
                                const GEOSPreparedGeometry *pg1, double x,
                                double y);

#endif

#endif