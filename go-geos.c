#include "go-geos.h"

// Using cgo to call C functions from Go has a high overhead. The functions in
// this file batch multiple calls to GEOS in C (rather than Go) to increase
// performance.

enum { bounds_MinX, bounds_MinY, bounds_MaxX, bounds_MaxY };

uintptr_t c_GEOSGeom_getUserData_r(GEOSContextHandle_t handle,
                                   const GEOSGeometry *g) {
  void *userdata = GEOSGeom_getUserData_r(handle, g);
  return (uintptr_t)userdata;
}

void c_GEOSGeom_setUserData_r(GEOSContextHandle_t handle, GEOSGeometry *g,
                              uintptr_t userdata) {
  GEOSGeom_setUserData_r(handle, g, (void *)userdata);
}

// c_GEOSGeomBounds_r extends bounds to include g.
void c_GEOSGeomBounds_r(GEOSContextHandle_t handle, const GEOSGeometry *g,
                        double *minX, double *minY, double *maxX,
                        double *maxY) {
  if (GEOSisEmpty_r(handle, g)) {
    return;
  }

  switch (GEOSGeomTypeId_r(handle, g)) {
  case GEOS_POINT: {
    double x;
    GEOSGeomGetX_r(handle, g, &x);
    if (x < *minX) {
      *minX = x;
    }
    if (x > *maxX) {
      *maxX = x;
    }
    double y;
    GEOSGeomGetY_r(handle, g, &y);
    if (y < *minY) {
      *minY = y;
    }
    if (y > *maxY) {
      *maxY = y;
    }
  } break;
  case GEOS_LINESTRING:
    // fallthrough
  case GEOS_LINEARRING: {
    const GEOSCoordSequence *s = GEOSGeom_getCoordSeq_r(handle, g);
    unsigned int size;
    GEOSCoordSeq_getSize_r(handle, s, &size);
    for (int i = 0; i < size; ++i) {
      double x;
      GEOSCoordSeq_getX_r(handle, s, i, &x);
      if (x < *minX) {
        *minX = x;
      }
      if (x > *maxX) {
        *maxX = x;
      }
      double y;
      GEOSCoordSeq_getY_r(handle, s, i, &y);
      if (y < *minY) {
        *minY = y;
      }
      if (y > *maxY) {
        *maxY = y;
      }
    }
  } break;
  case GEOS_POLYGON:
    c_GEOSGeomBounds_r(handle, GEOSGetExteriorRing_r(handle, g), minX, minY,
                       maxX, maxY);
    for (int i = 0, n = GEOSGetNumInteriorRings_r(handle, g); i < n; ++i) {
      c_GEOSGeomBounds_r(handle, GEOSGetInteriorRingN_r(handle, g, i), minX,
                         minY, maxX, maxY);
    }
    break;
  case GEOS_MULTIPOINT:
    // fallthrough
  case GEOS_MULTILINESTRING:
    // fallthrough
  case GEOS_MULTIPOLYGON:
    // fallthrough
  case GEOS_GEOMETRYCOLLECTION:
    for (int i = 0, n = GEOSGetNumGeometries_r(handle, g); i < n; ++i) {
      c_GEOSGeomBounds_r(handle, GEOSGetGeometryN_r(handle, g, i), minX, minY,
                         maxX, maxY);
    }
    break;
  }
}

// c_GEOSGeomGetInfo_r returns information about g. It returns 0 on any
// exception, 1 otherwise.
int c_GEOSGeomGetInfo_r(GEOSContextHandle_t handle, const GEOSGeometry *g,
                        int *typeID, int *numGeometries, int *numPoints,
                        int *numInteriorRings) {
  *typeID = GEOSGeomTypeId_r(handle, g);
  if (*typeID == -1) {
    return 0;
  }
  *numGeometries = GEOSGetNumGeometries_r(handle, g);
  if (*numGeometries == -1) {
    return 0;
  }
  switch (*typeID) {
  case GEOS_LINESTRING:
    // fallthrough
  case GEOS_LINEARRING:
    *numPoints = GEOSGeomGetNumPoints_r(handle, g);
    if (*numPoints == -1) {
      return 0;
    }
    break;
  case GEOS_POLYGON:
    *numInteriorRings = GEOSGetNumInteriorRings_r(handle, g);
    if (*numInteriorRings == -1) {
      return 0;
    }
    break;
  }
  return 1;
}

void c_errorMessageHandler(const char *message, void *userdata) {
  void go_errorMessageHandler(const char *, void *);
  go_errorMessageHandler(message, userdata);
}

// c_newGEOSGeomFromBounds_r returns a new GEOSGeom representing bounds. It
// returns NULL on any exception.
GEOSGeometry *c_newGEOSGeomFromBounds_r(GEOSContextHandle_t handle, int *typeID,
                                        double minX, double minY, double maxX,
                                        double maxY) {
  if (minX > maxX || minY > maxY) {
    *typeID = GEOS_POINT;
    return GEOSGeom_createEmptyPoint_r(handle);
  }
  if (minX == maxX && minY == maxY) {
    GEOSCoordSequence *s = GEOSCoordSeq_create_r(handle, 1, 2);
    if (s == NULL) {
      return NULL;
    }
    if (GEOSCoordSeq_setX_r(handle, s, 0, minX) == 0 ||
        GEOSCoordSeq_setY_r(handle, s, 0, minY) == 0) {
      GEOSCoordSeq_destroy_r(handle, s);
      return NULL;
    }
    GEOSGeometry *g = GEOSGeom_createPoint_r(handle, s);
    if (g == NULL) {
      GEOSCoordSeq_destroy_r(handle, s);
      return NULL;
    }
    *typeID = GEOS_POINT;
    return g;
  }
  const double flatCoords[10] = {minX, minY, maxX, minY, maxX,
                                 maxY, minX, maxY, minX, minY};
  GEOSCoordSequence *s =
      GEOSCoordSeq_copyFromBuffer_r(handle, flatCoords, 5, 0, 0);
  if (s == NULL) {
    return NULL;
  }
  GEOSGeometry *shell = GEOSGeom_createLinearRing_r(handle, s);
  if (shell == NULL) {
    GEOSCoordSeq_destroy_r(handle, s);
    return NULL;
  }
  GEOSGeometry *polygon = GEOSGeom_createPolygon_r(handle, shell, NULL, 0);
  if (polygon == NULL) {
    GEOSGeom_destroy_r(handle, shell);
    return NULL;
  }
  *typeID = GEOS_POLYGON;
  return polygon;
}

void c_GEOSSTRtree_query_callback(void *elem, void *userdata) {
  void go_GEOSSTRtree_query_callback(void *, void *);
  go_GEOSSTRtree_query_callback(elem, userdata);
}

int c_GEOSSTRtree_distance_callback(const void *item1, const void *item2,
                                    double *distance, void *userdata) {
  int go_GEOSSTRtree_distance_callback(const void *, const void *, double *,
                                       void *);
  return go_GEOSSTRtree_distance_callback(item1, item2, distance, userdata);
}

GEOSGeometry *c_GEOSMakeValidWithParams_r(GEOSContextHandle_t handle,
                                          const GEOSGeometry *g,
                                          enum GEOSMakeValidMethods method,
                                          int keepCollapsed) {
  GEOSGeometry *res;
  GEOSMakeValidParams *par;

  par = GEOSMakeValidParams_create_r(handle);
  GEOSMakeValidParams_setKeepCollapsed_r(handle, par, keepCollapsed);
  GEOSMakeValidParams_setMethod_r(handle, par, method);

  res = GEOSMakeValidWithParams_r(handle, g, par);

  GEOSMakeValidParams_destroy_r(handle, par);

  return res;
}

#if GEOS_VERSION_MAJOR < 3 ||                                                  \
    (GEOS_VERSION_MAJOR == 3 && GEOS_VERSION_MINOR < 11)

GEOSGeometry *GEOSConcaveHull_r(GEOSContextHandle_t handle,
                                const GEOSGeometry *g, double ratio,
                                unsigned int allowHoles) {
  return NULL;
}

#endif

#if GEOS_VERSION_MAJOR < 3 ||                                                  \
    (GEOS_VERSION_MAJOR == 3 && GEOS_VERSION_MINOR < 12)

GEOSGeometry *GEOSConcaveHullByLength_r(GEOSContextHandle_t handle,
                                        const GEOSGeometry *g, double ratio,
                                        unsigned int allowHoles) {
  return NULL;
}

char GEOSPreparedContainsXY_r(GEOSContextHandle_t handle,
                              const GEOSPreparedGeometry *pg1, double x,
                              double y) {
  return 0;
}

char GEOSPreparedIntersectsXY_r(GEOSContextHandle_t handle,
                                const GEOSPreparedGeometry *pg1, double x,
                                double y) {
  return 0;
}

#endif