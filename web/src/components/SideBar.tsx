import "./SideBar.scss"
import * as React from "react"
import { ListGroup, Badge, DropdownButton, ButtonGroup, Dropdown } from "react-bootstrap"
import { connect } from "react-redux"

import { bindActionCreators, Dispatch, AnyAction } from "redux"
import { StoreState, AggregatedError, SortFilters, ErrorsState, SeverityFilter } from "data/types"
import { setErrorsSortFilter, setErrorsSeverityFilter, setErrorsSearchFilter } from "data/errors"

interface DispatchProps {
  setErrorsSortFilter: (filter: SortFilters) => void
  setErrorsSeverityFilter: (severity: SeverityFilter) => void
  setErrorsSearchFilter: (searchTerm: string) => void
}

interface ConnectedProps {
  activeError: AggregatedError
  activeSortFilter: SortFilters
  activeSeverityFilter: ErrorsState["severityFilter"]
  searchKey: string
}

interface SidebarProps {
  errors: AggregatedError[]
  handleErrorSelect: (errorKey: string) => void
}

type Props = ConnectedProps & DispatchProps & SidebarProps

export const SORT_FILTERS = {
  "latest_occurrence": "Latest Occurence",
  "event_count": "Event Count",
}


const SideBar: React.FC<Props> = (props) => {

  const onSearchByKeyChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const { value } = event.target
    props.setErrorsSearchFilter(value)
  }

  const renderErrors = () => {
    if (!Array.isArray(props.errors)) {
      return <div>loading...</div>
    }

    if (props.errors.length === 0) {
      return <div>no errors returned from api</div>
    }

    return props.errors.map((error, index) => {
      return (
        <ListGroup.Item
          action
          className={`sidebar-item sidebar-item-${error.severity}`}
          onClick={_ => props.handleErrorSelect(error.aggregation_key)}
          active={props.activeError === undefined ? false : error.aggregation_key === props.activeError.aggregation_key}
          key={index}
        >
          {error.aggregation_key} <Badge variant="secondary" className="float-right">{error.total_count}</Badge>
        </ListGroup.Item>
      )
    })
  }

  const renderActions = () => {
    return (
      <div className="grid-component-actions">
        <div className="filters">
          <DropdownButton
            id="sort-btn"
            as={ButtonGroup}
            variant="secondary"
            title={`Sort by: ${SORT_FILTERS[props.activeSortFilter]}`}
            size="sm"
          >
            {Object.keys(SORT_FILTERS).map((filter: SortFilters) => (
              <Dropdown.Item
                key={`${filter}-sortFilter`}
                active={filter === props.activeSortFilter}
                onClick={_ => props.setErrorsSortFilter(filter)}
              >
                {SORT_FILTERS[filter]}
              </Dropdown.Item>
            ))}
          </DropdownButton>

          <DropdownButton
            id="severity-filter"
            as={ButtonGroup}
            variant="secondary"
            title={`Severity: ${props.activeSeverityFilter}`}
            size="sm"
          >
            {Object.keys(SeverityFilter).map((severitiy: SeverityFilter) => (
              <Dropdown.Item
                key={`${severitiy}-sortFilter`}
                active={severitiy === props.activeSeverityFilter}
                onClick={_ => props.setErrorsSeverityFilter(severitiy)}
              >
                {severitiy}
              </Dropdown.Item>
            ))}
          </DropdownButton>
        </div>

        <input
          onChange={onSearchByKeyChange}
          placeholder="Search for an error"
          value={props.searchKey}
        />
      </div>
    )
  }

  return (
    <div className="grid-component">
      {renderActions()}

      <ListGroup>
        {renderErrors()}
      </ListGroup>
    </div>
  )
}

const mapDispatchToProps = (dispatch: Dispatch<AnyAction>): DispatchProps => {
  return bindActionCreators({ setErrorsSortFilter, setErrorsSeverityFilter, setErrorsSearchFilter }, dispatch)
}

const mapStateToProps = (state: StoreState) => {
  return {
    activeError: state.errorsReducer.activeError,
    activeSortFilter: state.errorsReducer.activeSortFilter,
    activeSeverityFilter: state.errorsReducer.severityFilter,
    searchKey: state.errorsReducer.searchTerm
  }
}

export default connect(mapStateToProps, mapDispatchToProps)(SideBar)
