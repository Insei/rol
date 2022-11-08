package tests

import (
	"github.com/google/uuid"
	"testing"
)

var (
	repoTesters = []IGenericRepositoryTester{
		NewGormSQLiteGenericRepositoryTester[uuid.UUID, TestEntityUUID](),
		NewGormSQLiteGenericRepositoryTester[int, TestEntityInt](),
	}
)

func Test_GenericRepo_CountZero(t *testing.T) {
	for _, tester := range repoTesters {
		tester.SetTesting(t)
		tester._0Count_Zero()
	}
}

func Test_GenericRepo_InsertAndDelete(t *testing.T) {
	for _, tester := range repoTesters {
		tester.SetTesting(t)
		tester._1InsertAndDelete()
	}
}

func Test_GenericRepo_DeleteAll(t *testing.T) {
	for _, tester := range repoTesters {
		tester.SetTesting(t)
		tester._2DeleteAll()
	}
}

func Test_GenericRepo_GetByID_Found(t *testing.T) {
	for _, tester := range repoTesters {
		tester.SetTesting(t)
		tester._3GetByID_Found()
	}
}

func Test_GenericRepo_GetByID_NotFound(t *testing.T) {
	for _, tester := range repoTesters {
		tester.SetTesting(t)
		tester._4GetByID_NotFound()
	}
}

func Test_GenericRepo_GetByIDExtended_Found(t *testing.T) {
	for _, tester := range repoTesters {
		tester.SetTesting(t)
		tester._5GetByIDExtended_Found()
	}
}

func Test_GenericRepo_GetByIDExtended_NotFound(t *testing.T) {
	for _, tester := range repoTesters {
		tester.SetTesting(t)
		tester._6GetByIDExtended_NotFound()
	}
}

func Test_GenericRepo_GetByIDExtended_WithEmptyConditions(t *testing.T) {
	for _, tester := range repoTesters {
		tester.SetTesting(t)
		tester._7GetByIDExtended_WithEmptyConditions()
	}
}

func Test_GenericRepo_GetByIDExtended_WithNotExistedFieldCondition(t *testing.T) {
	for _, tester := range repoTesters {
		tester.SetTesting(t)
		tester._8GetByIDExtended_WithNotExistedFieldCondition()
	}
}

func Test_GenericRepo_GetByIDExtended_WithExistedFieldCondition_Found(t *testing.T) {
	for _, tester := range repoTesters {
		tester.SetTesting(t)
		tester._9GetByIDExtended_WithExistedFieldCondition_Found()
	}
}

func Test_GenericRepo_GetByIDExtended_WithExistedFieldCondition_NotFound(t *testing.T) {
	for _, tester := range repoTesters {
		tester.SetTesting(t)
		tester._10GetByIDExtended_WithExistedFieldCondition_NotFound()
	}
}

func Test_GenericRepo_GetByIDExtended_WithMultipleExistedFieldCondition_Found(t *testing.T) {
	for _, tester := range repoTesters {
		tester.SetTesting(t)
		tester._11GetByIDExtended_WithMultipleExistedFieldCondition_Found()
	}
}

func Test_GenericRepo_GetByIDExtended_WithMultipleExistedFieldCondition_NotFound(t *testing.T) {
	for _, tester := range repoTesters {
		tester.SetTesting(t)
		tester._12GetByIDExtended_WithMultipleExistedFieldCondition_NotFound()
	}
}

func Test_GenericRepo_GetByIDExtended_WithORMultipleExistedFieldCondition_NotFound(t *testing.T) {
	for _, tester := range repoTesters {
		tester.SetTesting(t)
		tester._13GetByIDExtended_WithORMultipleExistedFieldCondition_NotFound()
	}
}

func Test_GenericRepo_GetByIDExtended_WithORMultipleExistedFieldCondition_Found(t *testing.T) {
	for _, tester := range repoTesters {
		tester.SetTesting(t)
		tester._14GetByIDExtended_WithORMultipleExistedFieldCondition_Found()
	}
}

func Test_GenericRepo_Delete(t *testing.T) {
	for _, tester := range repoTesters {
		tester.SetTesting(t)
		tester._15Delete()
	}
}

func Test_GenericRepo_IsExist_Exist(t *testing.T) {
	for _, tester := range repoTesters {
		tester.SetTesting(t)
		tester._16IsExist_Exist()
	}
}

func Test_GenericRepo_IsExist_NotExist(t *testing.T) {
	for _, tester := range repoTesters {
		tester.SetTesting(t)
		tester._17IsExist_NotExist()
	}
}

func Test_GenericRepo_IsExist_WithEmptyConditions(t *testing.T) {
	for _, tester := range repoTesters {
		tester.SetTesting(t)
		tester._18IsExist_WithEmptyConditions_Exist()
	}
}

func Test_GenericRepo_IsExist_WithNotExistedFieldCondition(t *testing.T) {
	for _, tester := range repoTesters {
		tester.SetTesting(t)
		tester._19IsExist_WithNotExistedFieldCondition_Error()
	}
}

func Test_GenericRepo_IsExist_WithExistedFieldCondition_Exist(t *testing.T) {
	for _, tester := range repoTesters {
		tester.SetTesting(t)
		tester._20IsExist_WithExistedFieldCondition_Exist()
	}
}

func Test_GenericRepo_IsExist_WithExistedFieldCondition_NotExist(t *testing.T) {
	for _, tester := range repoTesters {
		tester.SetTesting(t)
		tester._21IsExist_WithExistedFieldCondition_NotExist()
	}
}

func Test_GenericRepo_IsExist_WithMultipleExistedFieldCondition_Exist(t *testing.T) {
	for _, tester := range repoTesters {
		tester.SetTesting(t)
		tester._22IsExist_WithMultipleExistedFieldCondition_Exist()
	}
}

func Test_GenericRepo_IsExist_WithMultipleExistedFieldCondition_NotExist(t *testing.T) {
	for _, tester := range repoTesters {
		tester.SetTesting(t)
		tester._23IsExist_WithMultipleExistedFieldCondition_NotExist()
	}
}

func Test_GenericRepo_IsExist_WithORMultipleExistedFieldCondition_NotExist(t *testing.T) {
	for _, tester := range repoTesters {
		tester.SetTesting(t)
		tester._24IsExist_WithORMultipleExistedFieldCondition_NotExist()
	}
}

func Test_GenericRepo_IsExist_WithORMultipleExistedFieldCondition_Exist(t *testing.T) {
	for _, tester := range repoTesters {
		tester.SetTesting(t)
		tester._25IsExist_WithORMultipleExistedFieldCondition_Exist()
	}
}

func Test_GenericRepo_IsExist_WithExistedNullFieldCondition_Exist(t *testing.T) {
	for _, tester := range repoTesters {
		tester.SetTesting(t)
		tester._26IsExist_WithExistedNullFieldCondition_Exist()
	}
}

func Test_GenericRepo_IsExist_WithExistedNullFieldCondition_NotExist(t *testing.T) {
	for _, tester := range repoTesters {
		tester.SetTesting(t)
		tester._27IsExist_WithExistedNullFieldCondition_NotExist()
	}
}

func Test_GenericRepo_IsExist_WithExistedNullableFieldCondition_Exist(t *testing.T) {
	for _, tester := range repoTesters {
		tester.SetTesting(t)
		tester._28IsExist_WithExistedNullableFieldCondition_Exist()
	}
}

func Test_GenericRepo_IsExist_WithExistedNullableFieldCondition_NotExist(t *testing.T) {
	for _, tester := range repoTesters {
		tester.SetTesting(t)
		tester._29IsExist_WithExistedNullableFieldCondition_NotExist()
	}
}

func Test_GenericRepo_Count_All(t *testing.T) {
	for _, tester := range repoTesters {
		tester.SetTesting(t)
		tester._30Count_All()
	}
}

func Test_GenericRepo_Count_WithNotExistedFieldCondition_Error(t *testing.T) {
	for _, tester := range repoTesters {
		tester.SetTesting(t)
		tester._31Count_WithNotExistedFieldCondition_Error()
	}
}

func Test_GenericRepo_Count_ById_OneElem(t *testing.T) {
	for _, tester := range repoTesters {
		tester.SetTesting(t)
		tester._32Count_ById_OneElem()
	}
}

func Test_GenericRepo_Count_ById_TwoElems(t *testing.T) {
	for _, tester := range repoTesters {
		tester.SetTesting(t)
		tester._33Count_ById_TwoElems()
	}
}

func Test_GenericRepo_Count_ByIteratorRange(t *testing.T) {
	for _, tester := range repoTesters {
		tester.SetTesting(t)
		tester._34Count_ByIteratorRange()
	}
}

func Test_GenericRepo_Clean(t *testing.T) {
	for _, tester := range repoTesters {
		err := tester.Dispose()
		if err != nil {
			t.Error("Failed to clean repository tester stuff")
		}
	}
}
