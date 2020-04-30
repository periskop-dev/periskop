import * as RemoteData from "data/remote-data";
import { SORT_FILTERS } from "components/SideBar"

export type Headers = {
  [key: string]: any
}

export interface HttpContext {
  "request_method"?: string,
  "request_headers"?: Headers,
  "request_url"?: string
}

export interface ErrorInstance {
  "class"?: string,
  "message"?: string,
  "stacktrace"?: string[],
  "cause"?: ErrorInstance
}

export interface Error {
  "error"?: ErrorInstance,
  "timestamp"?: number,
  "severity"?: string,
  "http_context"?: HttpContext
}

export interface AggregatedError {
  "aggregation_key"?: string,
  "total_count"?: number,
  "severity"?: string,
  "latest_errors"?: Error[]
}

export type ServicesState = {
  services: RemoteData.RemoteData<any, string[]>
} 

export type SortFilters = keyof typeof SORT_FILTERS

export type ErrorsState = {
  errors: RemoteData.RemoteData<any, AggregatedError[]>
  activeError?: AggregatedError,
  updatedAt?: number,
  activeService?: string,
  latestExceptionIndex: number
  activeSortFilter: SortFilters,  
}

export type StoreState = {
  servicesReducer: ServicesState,
  errorsReducer: ErrorsState
}