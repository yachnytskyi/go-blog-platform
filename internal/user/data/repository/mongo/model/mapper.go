package model

import (
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/pkg/interfaces"
	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	common "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	model "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/data/repository/mongo"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

func UserRepositoryToUsersRepositoryMapper(usersRepository []UserRepository) UsersRepository {
	return NewUsersRepository(
		append([]UserRepository{}, usersRepository...),
	)
}

func UsersRepositoryToUsersMapper(usersRepository UsersRepository) userModel.Users {
	users := make([]userModel.User, len(usersRepository.Users))
	for index, userRepository := range usersRepository.Users {
		users[index] = UserRepositoryToUserMapper(userRepository)
	}

	return userModel.NewUsers(
		users,
		usersRepository.PaginationResponse,
	)
}

func UserRepositoryToUserMapper(userRepository UserRepository) userModel.User {
	return userModel.NewUser(
		userRepository.ID.Hex(),
		userRepository.Username,
		userRepository.Email,
		userRepository.Password,
		userRepository.Role,
		userRepository.Verified,
		userRepository.CreatedAt,
		userRepository.UpdatedAt,
	)
}

func UserResetExpiryRepositoryToUserResetExpiryMapper(userResetExpiryRepository UserResetExpiryRepository) userModel.UserResetExpiry {
	return userModel.NewUserResetExpiry(
		userResetExpiryRepository.ResetExpiry,
	)
}

func UserCreateToUserCreateRepositoryMapper(userCreate userModel.UserCreate) UserCreateRepository {
	return NewUserCreateRepository(
		userCreate.Username,
		userCreate.Email,
		userCreate.Password,
		userCreate.Role,
		userCreate.Verified,
		userCreate.VerificationCode,
	)
}

func UserUpdateToUserUpdateRepositoryMapper(logger interfaces.Logger, location string, userUpdate userModel.UserUpdate) common.Result[UserUpdateRepository] {
	userObjectID := model.HexToObjectIDMapper(logger, location+".UserUpdateToUserUpdateRepositoryMapper", userUpdate.ID)
	if validator.IsError(userObjectID.Error) {
		return common.NewResultOnFailure[UserUpdateRepository](userObjectID.Error)
	}

	return common.NewResultOnSuccess(NewUserUpdateRepository(
		userObjectID.Data,
		userUpdate.Username,
	))
}

func UserForgottenPasswordToUserForgottenPasswordRepositoryMapper(userForgottenPassword userModel.UserForgottenPassword) UserForgottenPasswordRepository {
	return NewUserForgottenPasswordRepository(
		userForgottenPassword.ResetToken,
		userForgottenPassword.ResetExpiry,
	)
}

func UserResetPasswordToUserResetPasswordRepositoryMapper(userResetPassword userModel.UserResetPassword) UserResetPasswordRepository {
	return NewUserResetPasswordRepository(
		userResetPassword.Password,
	)
}
