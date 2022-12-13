#include "geos.h"

// Using cgo to call C functions from Go has a high overhead. The functions in
// this file batch multiple calls to GEOS in C (rather than Go) to increase
// performance.

enum { bounds_MinX, bounds_MinY, bounds_MaxX, bounds_MaxY };

// c_GEOSCoordSeq_getFlatCoords_r writes s's coordinate data to flatCoords,
// which must contain size*dims elements. It returns 0 on any exception, 1
// otherwise.
int c_GEOSCoordSeq_getFlatCoords_r(GEOSContextHandle_t handle,
                                   const GEOSCoordSequence *s,
                                   unsigned int size, unsigned int dims,
                                   double *flatCoords) {
#if GEOS_VERSION_MAJOR < 3 ||                                                  \
    (GEOS_VERSION_MAJOR == 3 && GEOS_VERSION_MINOR < 10)
  double *val = flatCoords;
  for (unsigned int idx = 0; idx < size; ++idx) {
    for (unsigned int dim = 0; dim < dims; ++dim) {
      if (GEOSCoordSeq_getOrdinate_r(handle, s, idx, dim, val++) == 0) {
        return 0;
      }
    }
  }
  return 1;
#else
  int hasZ = dims > 2;
  int hasM = dims > 3;
  return GEOSCoordSeq_copyToBuffer_r(handle, s, flatCoords, hasZ, hasM);
#endif
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

// c_newGEOSCoordSeqFromFlatCoords returns a new GEOSCoordSequence populated
// with flatCoords. It returns NULL on any exception.
GEOSCoordSequence *c_newGEOSCoordSeqFromFlatCoords_r(GEOSContextHandle_t handle,
                                                     unsigned int size,
                                                     unsigned int dims,
                                                     const double *flatCoords) {
#if GEOS_VERSION_MAJOR < 3 ||                                                  \
    (GEOS_VERSION_MAJOR == 3 && GEOS_VERSION_MINOR < 10)
  GEOSCoordSequence *s = GEOSCoordSeq_create_r(handle, size, dims);
  if (s == NULL) {
    return NULL;
  }
  const double *val = flatCoords;
  for (unsigned int idx = 0; idx < size; ++idx) {
    for (unsigned int dim = 0; dim < dims; ++dim) {
      if (GEOSCoordSeq_setOrdinate_r(handle, s, idx, dim, *val++) == 0) {
        GEOSCoordSeq_destroy_r(handle, s);
        return NULL;
      }
    }
  }
  return s;
#else
  int hasZ = dims > 2;
  int hasM = dims > 3;
  return GEOSCoordSeq_copyFromBuffer_r(handle, flatCoords, size, hasZ, hasM);
#endif
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
#if GEOS_VERSION_MAJOR < 3 ||                                                  \
    (GEOS_VERSION_MAJOR == 3 && GEOS_VERSION_MINOR < 10)
  GEOSCoordSequence *s = GEOSCoordSeq_create_r(handle, 5, 2);
  if (s == NULL) {
    return NULL;
  }
  const double *val = flatCoords;
  for (unsigned idx = 0; idx < 5; idx++) {
    if (GEOSCoordSeq_setX_r(handle, s, idx, *val++) == 0 ||
        GEOSCoordSeq_setY_r(handle, s, idx, *val++) == 0) {
      GEOSCoordSeq_destroy_r(handle, s);
      return NULL;
    }
  }
#else
  GEOSCoordSequence *s =
      GEOSCoordSeq_copyFromBuffer_r(handle, flatCoords, 5, 0, 0);
  if (s == NULL) {
    return NULL;
  }
#endif
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

#if GEOS_VERSION_MAJOR < 3 ||                                                  \
    (GEOS_VERSION_MAJOR == 3 && GEOS_VERSION_MINOR < 10)

GEOSGeometry *GEOSDensify_r(GEOSContextHandle_t handle, const GEOSGeometry *g,
                            double tolerance) {
  return NULL;
}

GEOSGeometry *GEOSDifferencePrec_r(GEOSContextHandle_t handle,
                                   const GEOSGeometry *g1,
                                   const GEOSGeometry *g2, double gridSize) {
  return NULL;
}

char GEOSDistanceWithin_r(GEOSContextHandle_t handle, const GEOSGeometry *g1,
                          const GEOSGeometry *g2, double dist) {
  return 2;
}

int GEOSFrechetDistance_r(GEOSContextHandle_t handle, const GEOSGeometry *g1,
                          const GEOSGeometry *g2, double *dist) {
  return 0;
}

int GEOSFrechetDistanceDensify_r(GEOSContextHandle_t handle,
                                 const GEOSGeometry *g1, const GEOSGeometry *g2,
                                 double densifyFrac, double *dist) {
  return 0;
}

GEOSGeometry *GEOSIntersectionPrec_r(GEOSContextHandle_t handle,
                                     const GEOSGeometry *g1,
                                     const GEOSGeometry *g2, double gridSize) {
  return NULL;
}

GEOSGeometry *GEOSMaximumInscribedCircle_r(GEOSContextHandle_t handle,
                                           const GEOSGeometry *g,
                                           double tolerance) {
  return NULL;
}

GEOSGeoJSONReader *GEOSGeoJSONReader_create_r(GEOSContextHandle_t handle) {
  return NULL;
}

void GEOSGeoJSONReader_destroy_r(GEOSContextHandle_t handle,
                                 GEOSGeoJSONReader *reader) {}

GEOSGeometry *GEOSGeoJSONReader_readGeometry_r(GEOSContextHandle_t handle,
                                               GEOSGeoJSONReader *reader,
                                               const char *geojson) {
  return NULL;
}

GEOSGeoJSONWriter *GEOSGeoJSONWriter_create_r(GEOSContextHandle_t handle) {
  return NULL;
}

void GEOSGeoJSONWriter_destroy_r(GEOSContextHandle_t handle,
                                 GEOSGeoJSONWriter *reader) {}

char *GEOSGeoJSONWriter_writeGeometry_r(GEOSContextHandle_t handle,
                                        GEOSGeoJSONWriter *writer,
                                        const GEOSGeometry *g, int indent) {
  return NULL;
}

#endif

#if GEOS_VERSION_MAJOR < 3 ||                                                  \
    (GEOS_VERSION_MAJOR == 3 && GEOS_VERSION_MINOR < 11)

GEOSGeometry *GEOSConcaveHull_r(GEOSContextHandle_t handle,
                                const GEOSGeometry *g, double ratio,
                                unsigned int allowHoles) {
  return NULL;
}

#endif