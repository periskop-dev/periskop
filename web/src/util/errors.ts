import { AggregatedError } from "data/types"

export const errorSortByLatestOccurrence = (errors: AggregatedError[]) => {
  return [...errors].sort((a, b) => b.latest_errors[0].timestamp - a.latest_errors[0].timestamp)
}

export const errorSortByEventCount = (errors: AggregatedError[]) => {
  return [...errors].sort((a, b) => a.total_count >= b.total_count ? -1 : 1)
}