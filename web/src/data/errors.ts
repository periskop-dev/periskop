import { Dispatch } from "redux";
import * as RemoteData from "data/remote-data";
import { registerReducer } from "data/store"
import { AggregatedError, ErrorsState, SeverityFilter, SortFilters } from "data/types"
import { errorSortByLatestOccurrence, errorSortByEventCount } from "util/errors"
import { ThunkDispatch } from "redux-thunk";
const METADATA = require('config/metadata.js');

export const FETCH = "periskop/errors/FETCH"
export const FETCH_SUCCESS = "periskop/errors/FETCH_SUCCESS"
export const FETCH_FAILURE = "periskop/errors/FETCH_FAILURE"
export const SET_ACTIVE_ERROR = "periskop/errors/SET_ACTIVE_ERROR"
export const SET_CURRENT_EXCEPTION_INDEX = "periskop/errors/SET_CURRENT_EXCEPTION_INDEX"
export const SET_ERRORS_SORT_FILTER = "periskop/errors/SET_ERRORS_SORT_FILTER"
export const SET_ERRORS_SEVERITY_FILTER = "periskop/errors/SET_ERRORS_SEVERITY_FILTER"
export const SET_ERRORS_SEARCH_FILTER = "periskop/errors/SET_ERRORS_SEARCH_FILTER"
export const RESOLVE_ERROR = "periskop/errors/RESOLVE_ERROR"
export const RESOLVE_ERROR_FAILURE = "periskop/errors/RESOLVE_ERROR_FAILURE"

export type ErrorsAction =
  | { type: typeof FETCH; service: string }
  | { type: typeof FETCH_SUCCESS; errors: AggregatedError[] }
  | { type: typeof FETCH_FAILURE; error: any }
  | { type: typeof SET_ACTIVE_ERROR; errorKey: string }
  | { type: typeof SET_CURRENT_EXCEPTION_INDEX; index: number }
  | { type: typeof SET_ERRORS_SORT_FILTER; filter: SortFilters }
  | { type: typeof SET_ERRORS_SEVERITY_FILTER; severity: ErrorsState["severityFilter"] }
  | { type: typeof SET_ERRORS_SEARCH_FILTER; searchTerm: string }
  | { type: typeof RESOLVE_ERROR; service: string, errorKey: String }
  | { type: typeof RESOLVE_ERROR_FAILURE; service: string, errorKey: String }

export const fetchErrors = (service: string) => {
  return(
    dispatch: Dispatch<ErrorsAction>
  ) => {
    dispatch(fetchingErrors(service))

    return fetch(`${parseHostName()}services/${service}/errors/`).then(response => {
      return response
        .json()
        .then(errors => dispatch(fetchedErrorsSuccessfully(errors, service)))
        .catch(err => dispatch(fetchedErrorsFailed(err)))
    })
  }
}

export const setActiveError = (errorKey: string) => {
  return { type: SET_ACTIVE_ERROR, errorKey }
}

export const setCurrentExceptionIndex = (index: number) => {
  return { type: SET_CURRENT_EXCEPTION_INDEX, index }
}

export function fetchingErrors(service: string): ErrorsAction {
  return { type: FETCH, service }
}

export function fetchedErrorsSuccessfully(errors: AggregatedError[], service: string): ErrorsAction {
  return { type: FETCH_SUCCESS, errors }
}

export function fetchedErrorsFailed(error: any): ErrorsAction {
  return { type: FETCH_FAILURE, error }
}

export const setErrorsSortFilter = (filter: SortFilters) => {
  return { type: SET_ERRORS_SORT_FILTER, filter }
}

export const setErrorsSeverityFilter = (severity: ErrorsState["severityFilter"]) => {
  return { type: SET_ERRORS_SEVERITY_FILTER, severity }
}

export const setErrorsSearchFilter = (searchTerm: ErrorsState["searchTerm"]) => {
  return { type: SET_ERRORS_SEARCH_FILTER, searchTerm }
}

export const resolveError = (service: string, errorKey: string) => {
  return (dispatch: ThunkDispatch<{}, {}, ErrorsAction>) => {
    dispatch({ type: RESOLVE_ERROR, service: service, errorKey: errorKey })

    return fetch(`${parseHostName()}services/${service}/errors/${errorKey}/`, {
      method: 'DELETE'
    })
      .then((_) => dispatch(fetchErrors(service)))
      .catch(err => dispatch({ type: RESOLVE_ERROR_FAILURE, service: service, errorKey: errorKey }))
  }
}

export const ErrorsSortActions = {
  "latest_occurrence": errorSortByLatestOccurrence,
  "event_count": errorSortByEventCount,
}

const initialState: ErrorsState = {
  errors: RemoteData.idle(),
  activeError: undefined,
  updatedAt: undefined,
  latestExceptionIndex: 0,
  activeSortFilter: "latest_occurrence",
  severityFilter: SeverityFilter.All,
  searchTerm: ""
}

function errorsReducer(state = initialState, action: ErrorsAction) {
  switch (action.type) {
    case FETCH:
      return {
        ...state,
        errors: RemoteData.load(),
        activeService: action.service
      }
    case FETCH_SUCCESS:
      const { activeError } = state
      const nextActiveError = action.errors.find(e => e.aggregation_key === activeError?.aggregation_key)

      return {
        ...state,
        errors: RemoteData.succeed(action.errors),
        activeError: nextActiveError,
        updatedAt: (new Date()).getTime()
      }
    case FETCH_FAILURE:
      return {
        errors: RemoteData.fail(action.error)
      }
    case SET_ACTIVE_ERROR:
      switch (state.errors.status) {
        case RemoteData.SUCCESS:
          return {
            ...state,
            activeError: state.errors.data.find(e => e.aggregation_key === action.errorKey),
            latestExceptionIndex: 0
          }
        default: {
          return state
        }
      }
    case SET_CURRENT_EXCEPTION_INDEX:
      return {
        ...state,
        latestExceptionIndex: action.index
      }
    case SET_ERRORS_SORT_FILTER:
      return {
        ...state,
        activeSortFilter: action.filter,
      }
    case SET_ERRORS_SEVERITY_FILTER:
      return {
        ...state,
        severityFilter: action.severity,
      }
    case SET_ERRORS_SEARCH_FILTER:
      return {
        ...state,
        searchTerm: action.searchTerm,
      }
    default:
      return state
  }
}

registerReducer("errorsReducer", errorsReducer)

const parseHostName = () => {
  let windowUrl = new URL(window.location.origin)
  if (windowUrl.hostname === METADATA.api_host ) {
    windowUrl.port = METADATA.api_port
  }

  return windowUrl
}
