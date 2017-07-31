package naprrql

import (
	"github.com/nsip/nias2/naprr"
	"github.com/nsip/nias2/xml"
	"github.com/playlyfe/go-graphql"
)

func buildResolvers() map[string]interface{} {

	resolvers := map[string]interface{}{}

	resolvers["NaplanData/score_summaries"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return getObjects(getIdentifiers("NAPTestScoreSummary"))
	}

	resolvers["NaplanData/score_summaries_count"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return len(getIdentifiers("NAPTestScoreSummary")), nil
	}

	resolvers["NaplanData/students"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return getObjects(getIdentifiers("StudentPersonal"))
	}

	resolvers["NaplanData/students_count"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return len(getIdentifiers("StudentPersonal")), nil
	}

	resolvers["NaplanData/events"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return getObjects(getIdentifiers("NAPEventStudentLink"))
	}

	resolvers["NaplanData/events_count"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return len(getIdentifiers("NAPEventStudentLink")), nil
	}

	resolvers["NaplanData/responses"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return getObjects(getIdentifiers("NAPStudentResponseSet"))
	}

	resolvers["NaplanData/responses_count"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return len(getIdentifiers("NAPStudentResponseSet")), nil
	}

	resolvers["NaplanData/tests_count"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return len(getIdentifiers("NAPTest")), nil
	}

	resolvers["NaplanData/tests"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return getObjects(getIdentifiers("NAPTest"))
	}

	resolvers["NaplanData/schools_count"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return len(getIdentifiers("SchoolInfo")), nil
	}

	resolvers["NaplanData/schools"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return getObjects(getIdentifiers("SchoolInfo"))
	}

	resolvers["NaplanData/testlets_count"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return len(getIdentifiers("NAPTestlet")), nil
	}

	resolvers["NaplanData/testlets"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return getObjects(getIdentifiers("NAPTestlet"))
	}

	resolvers["NaplanData/testitems_count"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return len(getIdentifiers("NAPTestItem")), nil
	}

	resolvers["NaplanData/testitems"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return getObjects(getIdentifiers("NAPTestItem"))
	}

	resolvers["NaplanData/codeframes_count"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return len(getIdentifiers("NAPCodeFrame")), nil
	}

	resolvers["NaplanData/codeframes"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return getObjects(getIdentifiers("NAPCodeFrame"))
	}

	//
	// addition to spec that allows the original Item to be available when
	// reviewing item responses, e.g. to compare item correct response, item type etc.
	//
	resolvers["NAPResponseSet_ItemResponse/Item"] = func(params *graphql.ResolveParams) (interface{}, error) {

		linkedItem := make([]string, 0)
		// log.Printf("params: %#v\n\n", params)
		if napResponse, ok := params.Source.(xml.NAPResponseSet_ItemResponse); ok {
			linkedItem = append(linkedItem, napResponse.ItemRefID)
			obj, err := getObjects(linkedItem)
			return obj[0], err
		}
		return linkedItem, nil

	}

	//
	// shorthand lookup objects for basic school info
	//
	resolvers["NaplanData/school_details"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return getObjects(getIdentifiers("SchoolDetails"))
	}

	//
	// resolver for score summary report object
	//
	resolvers["NaplanData/score_summary_report_by_school"] = func(params *graphql.ResolveParams) (interface{}, error) {

		// get the acara ids from the request params
		acaraids := make([]string, 0)
		for _, a_id := range params.Args["acaraIDs"].([]interface{}) {
			acaraid, _ := a_id.(string)
			acaraids = append(acaraids, acaraid)
		}

		// get the sif refid for each of the acarids supplied
		refids := make([]string, 0)
		for _, acaraid := range acaraids {
			refid := getIdentifiers(acaraid)[0]
			refids = append(refids, refid)
		}

		// now construct the composite keys
		school_summary_keys := make([]string, 0)
		for _, refid := range refids {
			school_summary_keys = append(school_summary_keys, refid+":NAPTestScoreSummary")
		}

		summ_refids := make([]string, 0)
		for _, summary_key := range school_summary_keys {
			ids := getIdentifiers(summary_key)
			for _, id := range ids {
				summ_refids = append(summ_refids, id)
			}
		}

		summaries, err := getObjects(summ_refids)
		summary_datasets := make([]naprr.ScoreSummaryDataSet, 0)
		for _, summary := range summaries {
			summ, _ := summary.(xml.NAPTestScoreSummary)
			testid := []string{summ.NAPTestRefId}
			obj, _ := getObjects(testid)
			test, _ := obj[0].(xml.NAPTest)
			sds := naprr.ScoreSummaryDataSet{Summ: summ, Test: test}
			summary_datasets = append(summary_datasets, sds)
		}

		return summary_datasets, err

	}

	resolvers["NaplanData/school_infos_by_acaraid"] = func(params *graphql.ResolveParams) (interface{}, error) {

		// get the acara ids from the request params
		acaraids := make([]string, 0)
		for _, a_id := range params.Args["acaraIDs"].([]interface{}) {
			acaraid, _ := a_id.(string)
			acaraids = append(acaraids, acaraid)
		}

		// get the sif refid for each of the acarids supplied
		refids := make([]string, 0)
		for _, acaraid := range acaraids {
			refid := getIdentifiers(acaraid)[0]
			refids = append(refids, refid)
		}

		// get the school infos from the datastore
		siObjects, err := getObjects(refids)
		schoolInfos := make([]xml.SchoolInfo, 0)
		for _, sio := range siObjects {
			si, _ := sio.(xml.SchoolInfo)
			schoolInfos = append(schoolInfos, si)
		}

		return schoolInfos, err

	}

	resolvers["NaplanData/students_by_school"] = func(params *graphql.ResolveParams) (interface{}, error) {
		// get the acara ids from the request params
		acaraids := make([]string, 0)
		for _, a_id := range params.Args["acaraIDs"].([]interface{}) {
			acaraid, _ := a_id.(string)
			acaraids = append(acaraids, acaraid)
		}

		// get students for the schools
		studentids := make([]string, 0)
		for _, acaraid := range acaraids {
			key := "student_by_acaraid:" + acaraid
			studentRefIds := getIdentifiers(key)
			for _, refid := range studentRefIds {
				studentids = append(studentids, refid)
			}
		}

		return getObjects(studentids)

	}

	resolvers["NaplanData/domain_scores_report_by_school"] = func(params *graphql.ResolveParams) (interface{}, error) {

		// get the acara ids from the request params
		acaraids := make([]string, 0)
		for _, a_id := range params.Args["acaraIDs"].([]interface{}) {
			acaraid, _ := a_id.(string)
			acaraids = append(acaraids, acaraid)
		}

		// get students for the schools
		studentids := make([]string, 0)
		for _, acaraid := range acaraids {
			key := "student_by_acaraid:" + acaraid
			studentRefIds := getIdentifiers(key)
			for _, refid := range studentRefIds {
				studentids = append(studentids, refid)
			}
		}

		// get responses for student
		responseids := make([]string, 0)
		for _, studentid := range studentids {
			key := "responseset_by_student:" + studentid
			responseRefId := getIdentifiers(key)[0]
			responseids = append(responseids, responseRefId)
		}

		// get responses
		responses, err := getObjects(responseids)
		if err != nil {
			return []interface{}{}, err
		}

		// construct RDS by including referenced test
		results := make([]naprr.ResponseDataSet, 0)
		for _, response := range responses {
			resp, _ := response.(xml.NAPResponseSet)
			tests, err := getObjects([]string{resp.TestID})
			test, ok := tests[0].(xml.NAPTest)
			if err != nil || !ok {
				return []interface{}{}, err
			}
			rds := naprr.ResponseDataSet{Test: test, Response: resp}
			results = append(results, rds)
		}

		return results, nil

	}

	return resolvers
}
