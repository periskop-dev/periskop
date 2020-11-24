import { ErrorsSortActions } from "data/errors"
import { AggregatedError, ErrorsState, SeverityFilter, SortFilters } from "data/types"

export const errorSortByLatestOccurrence = (errors: AggregatedError[]) => {
  return [...errors].sort((a, b) => b.latest_errors[0].timestamp - a.latest_errors[0].timestamp)
}

export const errorSortByEventCount = (errors: AggregatedError[]) => {
  return [...errors].sort((a, b) => a.total_count >= b.total_count ? -1 : 1)
}

export const filterErrorsBySubstringMatch = (errors: AggregatedError[], searchTerm: string) => {
  return errors.filter((error) => error.aggregation_key.toLowerCase().includes(searchTerm.toLowerCase()))
}

export const filterErrorsBySeverity = (errors: AggregatedError[], severity: ErrorsState["severityFilter"]) => {
  if (severity === SeverityFilter.All) return errors

  return errors.filter((error) => error.severity.toLowerCase() === severity.toLowerCase())
}

export const getFilteredErrors = (
  errors: AggregatedError[],
  searchTerm: string,
  severity: SeverityFilter,
  sortFilter: SortFilters,
) => {
  const sortedErrors = ErrorsSortActions[sortFilter](errors)
  const searchedErrors = filterErrorsBySubstringMatch(sortedErrors, searchTerm)
  const filteredErrors = filterErrorsBySeverity(searchedErrors, severity)
  return filteredErrors
}