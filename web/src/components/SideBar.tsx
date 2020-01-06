import "SideBar.scss"
import * as React from "react";
import { ListGroup, Badge } from "react-bootstrap";

import { bindActionCreators, Dispatch, AnyAction } from "redux"
import { connect } from "react-redux"
import { StoreState, AggregatedError } from "data/types"
import { setActiveError } from "data/errors"
import { RouteComponentProps, withRouter } from "react-router";

interface DispatchProps {
  setActiveError: (notifcation: string) => void
}

interface ConnectedProps {
  activeError: AggregatedError
}

interface DefaultsProps {
  errors: AggregatedError[],
  handleErrorSelect: (errorKey: string) => void
}

type Props = ConnectedProps & DispatchProps & DefaultsProps & RouteComponentProps<{service: string}>

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
    } else {
      return props.errors.map((error, index) => {
        return (
          <ListGroup.Item action className={"sidebar-item" + " " + sidebarItemClass(error)} onClick={_ => props.handleErrorSelect(error.aggregation_key)} active={ props.activeError === undefined ? false : error.aggregation_key === props.activeError.aggregation_key } key={index}>
            {error.aggregation_key} <Badge variant="secondary" className="float-right">{error.total_count}</Badge>
          </ListGroup.Item>
        )
      })
    }
  }

  return (
    <div className="grid-component">
      <ListGroup>
        {renderNavItems()}
      </ListGroup>
    </div>
  )
}

const matchDispatchToProps = (dispatch: Dispatch<AnyAction>): DispatchProps => {
  return bindActionCreators({ setActiveError }, dispatch)
}

const mapStateToProps = (state: StoreState) => {
  return {
    activeError: state.errorsReducer.activeError,
  }
}

export default withRouter(connect<ConnectedProps, {}, RouteComponentProps<{service: string}>>(mapStateToProps, matchDispatchToProps)(SideBar))
