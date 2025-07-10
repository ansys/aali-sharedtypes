.. _logging:

logging
=======

This package provides functionality for logging.

Functions
---------

.. function:: InitLogger

   func InitLogger(GlobalConfig \*config.Config)

   ///////////////////////////////// Create Logger ///////////////////////////////// InitLogger initializes the global logger.  The function creates a new zap logger with the specified configuration and sets the global logger variable to the new logger.  Parameters: - GlobalConfig: The global configuration from the config package.


Types
-----

.. type:: ContextKey

   type ContextKey string

   ContextKey defines the supported context keys.


.. type:: Config

   type Config struct

   Config represents the configuration for the logging package.


.. type:: ContextMap

   type ContextMap struct

   ContextMap represents a context for managing key-value pairs with specific context keys. It allows setting, retrieving, and copying context data associated with various keys.


.. type:: Point

   type Point struct

   Point represents a data point in a time series metric.


.. type:: Resource

   type Resource struct

   Resource represents a named resource associated with a metric.


.. type:: Metric

   type Metric struct

   Metric represents a time series metric.


.. type:: Metrics

   type Metrics struct

   Metrics represents a collection of metrics.


Variables
---------

.. data:: Log

   var Log loggerWrapper

   Initialize the global logger variable.


.. data:: ERROR_FILE_LOCATION

   var ERROR\_FILE\_LOCATION string

   Initialize config variables


.. data:: LOG_LEVEL

   var LOG\_LEVEL string


.. data:: LOCAL_LOGS

   var LOCAL\_LOGS bool


.. data:: LOCAL_LOGS_LOCATION

   var LOCAL\_LOGS\_LOCATION string


.. data:: DATADOG_LOGS

   var DATADOG\_LOGS bool


.. data:: DATADOG_SOURCE

   var DATADOG\_SOURCE string


.. data:: DATADOG_STAGE

   var DATADOG\_STAGE string


.. data:: DATADOG_VERSION

   var DATADOG\_VERSION string


.. data:: DATADOG_SERVICE_NAME

   var DATADOG\_SERVICE\_NAME string


.. data:: DATADOG_API_KEY

   var DATADOG\_API\_KEY string


.. data:: DATADOG_LOGS_URL

   var DATADOG\_LOGS\_URL string


.. data:: DATADOG_METRICS

   var DATADOG\_METRICS bool


.. data:: DATADOG_METRICS_URL

   var DATADOG\_METRICS\_URL string
