import * as RemoteData from "data/remote-data";

export type Headers = {
  [key: string]: any
}

export interface HttpContext {
  "request_method"?: string,
  "request_headers"?: Headers,
  "request_url"?: string
}

export interface Error {
  "error"?: {
    "class"?: string,
    "message"?: string,
    "stacktrace"?: string[],
    "cause"?: Error
  },
  "timestamp"?: number,
  "severity"?: string,
  "http_context"?: HttpContext
}

export interface AggregatedError {
  "aggregation_key"?: string,
  "total_count"?: number,
  "severity"?: string,
  "latest_errors"?: Error
}

export type ServicesState = {
  services: RemoteData.RemoteData<any, string[]>
} 

export type ErrorsState = {
  errors: RemoteData.RemoteData<any, AggregatedError[]>
  activeError?: AggregatedError,
  updatedAt?: number,
  activeService?: string
}

export type StoreState = {
  servicesReducer: ServicesState,
  errorsReducer: ErrorsState
}