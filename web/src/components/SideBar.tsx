import "SideBar.scss"
import * as React from "react";
import { ListGroup, Badge, DropdownButton, ButtonGroup, Dropdown } from "react-bootstrap";

import { bindActionCreators, Dispatch, AnyAction } from "redux"
import { connect } from "react-redux"
import { StoreState, AggregatedError, SortFilters } from "data/types"
import { setActiveError, setActiveErrorSortFilter } from "data/errors"
import { RouteComponentProps, withRouter } from "react-router";

interface DispatchProps {
  setActiveError: (notifcation: string) => void
  setActiveErrorSortFilter: (filter: SortFilters) => void,
}

interface ConnectedProps {
  activeError: AggregatedError
}

interface DefaultsProps {
  errors: AggregatedError[],
  activeSortFilter: SortFilters,
  setActiveErrorSortFilter: (filter: SortFilters) => void,
  handleErrorSelect: (errorKey: string) => void
}

type Props = ConnectedProps & DispatchProps & DefaultsProps & RouteComponentProps<{ service: string }>

export const SORT_FILTERS = {
  "latest_occurrence": "Latest Occurence",
  "event_count": "Event Count",
}

const sidebarItemClass = (error: AggregatedError): string => {
  if (error.severity === "info") {
    return "sidebar-item-info"
  } else if (error.severity === "warning") {
    return "sidebar-item-warning"
  } else {
    return "sidebar-item-error"
  }
}

const SideBar = (props: Props) => {

  const renderNavItems = () => {
    if (props.errors.length === 0) {
      return <div>no errors returned from api</div>
    }

    return props.errors.map((error, index) => {
      return (
        <ListGroup.Item
          action className={"sidebar-item" + " " + sidebarItemClass(error)}
          onClick={_ => props.handleErrorSelect(error.aggregation_key)}
          active={props.activeError === undefined ? false : error.aggregation_key === props.activeError.aggregation_key}
          key={index}
        >
          {error.aggregation_key} <Badge variant="secondary" className="float-right">{error.total_count}</Badge>
        </ListGroup.Item>
      )
    })
  }

  const renderActions = () =>  {
    return (
      <div className='grid-component-actions'>
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
              onClick={_ => props.setActiveErrorSortFilter(filter)}
            >
              {SORT_FILTERS[filter]}
            </Dropdown.Item>
          ))}
        </DropdownButton>
      </div>
      )
  }

  return (
    <div className="grid-component">
      {renderActions()}

      <ListGroup>
        {renderNavItems()}
      </ListGroup>
    </div>
  )
}

const matchDispatchToProps = (dispatch: Dispatch<AnyAction>): DispatchProps => {
  return bindActionCreators({ setActiveError, setActiveErrorSortFilter }, dispatch)
}

const mapStateToProps = (state: StoreState) => {
  return {
    activeError: state.errorsReducer.activeError,
    activeSortFilter: state.errorsReducer.activeSortFilter,
  }
}

export default withRouter(connect<ConnectedProps, {}, RouteComponentProps<{service: string}>>(mapStateToProps, matchDispatchToProps)(SideBar))
